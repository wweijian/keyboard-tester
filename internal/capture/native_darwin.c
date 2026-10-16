//go:build darwin && cgo

#include "native_darwin.h"
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/IOKitLib.h>
#include <IOKit/hid/IOHIDManager.h>
#include <IOKit/hidsystem/IOHIDLib.h>
#include <mach/mach_error.h>
#include <stdlib.h>
#include <string.h>

struct kt_context {
    IOHIDManagerRef manager;
    CFArrayRef devices;
    CFMutableArrayRef opened;
    kt_event queue[queue_capacity];
    int head;
    int count;
    int failure;
};

static bool is_keyboard(IOHIDDeviceRef device) {
    return IOHIDDeviceConformsTo(device, kHIDPage_GenericDesktop, kHIDUsage_GD_Keyboard);
}

static bool is_key_device(IOHIDDeviceRef device) {
    return is_keyboard(device) ||
        IOHIDDeviceConformsTo(device, kHIDPage_Consumer, kHIDUsage_Csmr_ConsumerControl);
}

static bool is_virtual(IOHIDDeviceRef device) {
    CFTypeRef transport = IOHIDDeviceGetProperty(device, CFSTR(kIOHIDTransportKey));
    return transport && CFEqual(transport, CFSTR(kIOHIDTransportVirtualValue));
}

// USB keyboard and consumer-control interfaces share their physical USB parent.
// Other transports use their HID service's unique registry identity.
static uint64_t physical_id(IOHIDDeviceRef device) {
    io_service_t service = IOHIDDeviceGetService(device);
    uint64_t fallback = 0;
    IORegistryEntryGetRegistryEntryID(service, &fallback);
    IOObjectRetain(service);
    for (int depth = 0; depth < 16; depth++) {
        if (IOObjectConformsTo(service, "IOUSBHostDevice") ||
            IOObjectConformsTo(service, "IOBluetoothHIDDriver")) {
            uint64_t identity = fallback;
            IORegistryEntryGetRegistryEntryID(service, &identity);
            IOObjectRelease(service);
            return identity;
        }
        io_registry_entry_t parent = 0;
        kern_return_t result = IORegistryEntryGetParentEntry(service, kIOServicePlane, &parent);
        IOObjectRelease(service);
        if (result != KERN_SUCCESS) return fallback;
        service = parent;
    }
    IOObjectRelease(service);
    return fallback;
}

static void string_property(IOHIDDeviceRef device, CFStringRef key, char *out, size_t length) {
    CFTypeRef value = IOHIDDeviceGetProperty(device, key);
    if (value && CFGetTypeID(value) == CFStringGetTypeID())
        CFStringGetCString(value, out, length, kCFStringEncodingUTF8);
}

static uint32_t number_property(IOHIDDeviceRef device, CFStringRef key) {
    CFTypeRef value = IOHIDDeviceGetProperty(device, key);
    int32_t result = 0;
    if (value && CFGetTypeID(value) == CFNumberGetTypeID())
        CFNumberGetValue(value, kCFNumberSInt32Type, &result);
    return (uint32_t) result;
}

kt_context *kt_create(void) {
    kt_context *ctx = calloc(1, sizeof(*ctx));
    if (!ctx) return NULL;
    ctx->manager = IOHIDManagerCreate(kCFAllocatorDefault, kIOHIDOptionsTypeNone);
    ctx->opened = CFArrayCreateMutable(NULL, 0, &kCFTypeArrayCallBacks);
    if (!ctx->manager || !ctx->opened) {
        kt_destroy(ctx);
        return NULL;
    }
    IOHIDManagerSetDeviceMatching(ctx->manager, NULL);
    IOHIDManagerScheduleWithRunLoop(ctx->manager, CFRunLoopGetCurrent(), kCFRunLoopDefaultMode);
    return ctx;
}

int kt_refresh(kt_context *ctx) {
    if (ctx->devices) CFRelease(ctx->devices);
    ctx->devices = NULL;
    CFSetRef devices = IOHIDManagerCopyDevices(ctx->manager);
    if (!devices) return 0;
    CFIndex count = CFSetGetCount(devices);
    const void **values = calloc(count ? count : 1, sizeof(*values));
    if (!values) { CFRelease(devices); return -1; }
    CFSetGetValues(devices, values);
    ctx->devices = CFArrayCreate(NULL, values, count, &kCFTypeArrayCallBacks);
    free(values);
    CFRelease(devices);
    if (!ctx->devices) return -1;
    return (int) count;
}

bool kt_get_device(kt_context *ctx, int index, kt_device *info) {
    IOHIDDeviceRef device = (IOHIDDeviceRef) CFArrayGetValueAtIndex(ctx->devices, index);
    if (!is_keyboard(device) || is_virtual(device)) return false;
    memset(info, 0, sizeof(*info));
    info->id = physical_id(device);
    info->vendor = number_property(device, CFSTR(kIOHIDVendorIDKey));
    info->product = number_property(device, CFSTR(kIOHIDProductIDKey));
    string_property(device, CFSTR(kIOHIDProductKey), info->name, sizeof(info->name));
    string_property(device, CFSTR(kIOHIDTransportKey), info->transport, sizeof(info->transport));
    return info->id != 0;
}

static void input_value(void *context, IOReturn result, void *sender, IOHIDValueRef value) {
    (void) sender;
    kt_context *ctx = context;
    if (result != kIOReturnSuccess || IOHIDValueGetLength(value) > 8) return;
    IOHIDElementRef element = IOHIDValueGetElement(value);
    uint32_t usage = IOHIDElementGetUsage(element);
    uint32_t page = IOHIDElementGetUsagePage(element);
    if (!usage || !page) return;
    if (ctx->count == queue_capacity) { ctx->failure = event_overflow; return; }
    int position = (ctx->head + ctx->count) % queue_capacity;
    ctx->queue[position] = (kt_event) {
        .page = page, .usage = usage, .down = IOHIDValueGetIntegerValue(value) != 0,
    };
    ctx->count++;
}

static void device_removed(void *context, IOReturn result, void *sender) {
    (void) result;
    (void) sender;
    ((kt_context *) context)->failure = event_removed;
}

static IOReturn open_device(kt_context *ctx, IOHIDDeviceRef device) {
    IOReturn result = IOHIDDeviceOpen(device, kIOHIDOptionsTypeSeizeDevice);
    if (result != kIOReturnSuccess) return result;
    CFArrayAppendValue(ctx->opened, device);
    IOHIDDeviceRegisterInputValueCallback(device, input_value, ctx);
    IOHIDDeviceRegisterRemovalCallback(device, device_removed, ctx);
    IOHIDDeviceScheduleWithRunLoop(device, CFRunLoopGetCurrent(), kCFRunLoopDefaultMode);
    return kIOReturnSuccess;
}

uint32_t kt_start(kt_context *ctx, uint64_t id) {
    kt_stop(ctx);
    if (!id) return kIOReturnBadArgument;
    if (IOHIDCheckAccess(kIOHIDRequestTypeListenEvent) != kIOHIDAccessTypeGranted &&
        !IOHIDRequestAccess(kIOHIDRequestTypeListenEvent)) return kIOReturnNotPermitted;
    if (!ctx->devices) return kIOReturnNoDevice;
    CFIndex count = CFArrayGetCount(ctx->devices);
    for (CFIndex i = 0; i < count; i++) {
        IOHIDDeviceRef device = (IOHIDDeviceRef) CFArrayGetValueAtIndex(ctx->devices, i);
        if (!is_key_device(device) || physical_id(device) != id) continue;
        IOReturn result = open_device(ctx, device);
        if (result != kIOReturnSuccess) { kt_stop(ctx); return result; }
    }
    if (!CFArrayGetCount(ctx->opened)) return kIOReturnNoDevice;
    return kIOReturnSuccess;
}

void kt_stop(kt_context *ctx) {
    if (!ctx->opened) return;
    for (CFIndex i = 0; i < CFArrayGetCount(ctx->opened); i++) {
        IOHIDDeviceRef device = (IOHIDDeviceRef) CFArrayGetValueAtIndex(ctx->opened, i);
        IOHIDDeviceUnscheduleFromRunLoop(device, CFRunLoopGetCurrent(), kCFRunLoopDefaultMode);
        IOHIDDeviceRegisterInputValueCallback(device, NULL, NULL);
        IOHIDDeviceRegisterRemovalCallback(device, NULL, NULL);
        IOHIDDeviceClose(device, kIOHIDOptionsTypeNone);
    }
    CFArrayRemoveAllValues(ctx->opened);
    ctx->head = ctx->count = ctx->failure = 0;
}

void kt_destroy(kt_context *ctx) {
    if (!ctx) return;
    kt_stop(ctx);
    if (ctx->opened) CFRelease(ctx->opened);
    if (ctx->devices) CFRelease(ctx->devices);
    if (ctx->manager) {
        IOHIDManagerUnscheduleFromRunLoop(ctx->manager, CFRunLoopGetCurrent(), kCFRunLoopDefaultMode);
        CFRelease(ctx->manager);
    }
    free(ctx);
}

void kt_pump(void) {
    CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0.001, true);
}

bool kt_next(kt_context *ctx, kt_event *event) {
    if (ctx->failure) {
        *event = (kt_event) { .kind = ctx->failure };
        kt_stop(ctx);
        return true;
    }
    if (!ctx->count) return false;
    *event = ctx->queue[ctx->head];
    ctx->head = (ctx->head + 1) % queue_capacity;
    ctx->count--;
    return true;
}

const char *kt_error(uint32_t code) {
    return mach_error_string((mach_error_t) code);
}
