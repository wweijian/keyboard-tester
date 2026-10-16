#ifndef KEYBOARD_TESTER_NATIVE_H
#define KEYBOARD_TESTER_NATIVE_H

#include <stdbool.h>
#include <stdint.h>

typedef struct kt_context kt_context;

enum { queue_capacity = 4096, event_removed = 1, event_overflow = 2 };

typedef struct {
    uint64_t id;
    uint32_t vendor;
    uint32_t product;
    char name[256];
    char transport[64];
} kt_device;

typedef struct {
    uint32_t page;
    uint32_t usage;
    bool down;
    int kind;
} kt_event;

kt_context *kt_create(void);
void kt_destroy(kt_context *ctx);
int kt_refresh(kt_context *ctx);
bool kt_get_device(kt_context *ctx, int index, kt_device *info);
uint32_t kt_start(kt_context *ctx, uint64_t id);
void kt_stop(kt_context *ctx);
void kt_pump(void);
bool kt_next(kt_context *ctx, kt_event *event);
const char *kt_error(uint32_t code);

#endif
