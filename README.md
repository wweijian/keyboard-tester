# Keyboard Lab

A macOS keyboard tester written in Go, using Bubble Tea v2 and Lip Gloss. Select a physical keyboard, press its keys, and see which keys have been tested. Unknown keys and keys outside the selected preset appear with their raw HID identifiers.

## Build and run

Requires macOS 11 or later, Go 1.25 or later, and Xcode or the Xcode Command Line Tools. The native capture bridge requires cgo, which is enabled by default for local macOS builds.

```sh
make build
sudo ./bin/keyboard-tester --keyboard 65
```

Supported sizes: `40`, `60`, `65`, `75`, `80`, `96`, `100`. A trailing `%` is accepted, for example `--keyboard 100%`.

```sh
./bin/keyboard-tester          # help; does not access keyboards
./bin/keyboard-tester --help
./bin/keyboard-tester --list   # enumerate devices; does not start capture
```

In the UI, use the arrow keys to choose a device and Enter or the **Start** button to begin. Only the selected keyboard contributes results. A `✓` and mint colour mark a tested key; `·` marks a pending key. Shifted and unshifted characters share the same physical key. Left/right modifiers and keypad keys are distinct.

## macOS permissions and capture

Interactive testing requires **administrator privileges** because macOS restricts seizing a keyboard to privileged clients, as shown in [Apple's HID implementation](https://github.com/apple-oss-distributions/IOHIDFamily/blob/main/IOHIDFamily/IOHIDLibUserClient.cpp#L656). Launch the built binary with `sudo`; running as an ordinary user fails before opening the UI and prints the required command. Help and `--list` work without `sudo`. Build normally with `make build`; do not run the build or `go run` with `sudo`.

On the first capture attempt, macOS may also request **Input Monitoring** permission. Grant it to the terminal application or tester binary under **System Settings → Privacy & Security → Input Monitoring**, then quit and relaunch the tester with `sudo`. The UI distinguishes missing administrator privileges, denied Input Monitoring, busy devices, and unavailable devices.

Capture uses IOKit's [exclusive device access](https://developer.apple.com/documentation/iokit/1588670-iohiddeviceopen), which [prevents the system and other clients from receiving the device's events](https://developer.apple.com/documentation/kernel/1644502-anonymous/kiohidoptionstypeseizedevice). This suppresses the selected keyboard's shortcuts while testing without editing saved settings. Capture is system-wide while active: release it with **Devices** before switching to another application. The tester never falls back to unsuppressed capture if exclusive access fails.

USB keyboard interfaces and companion consumer-control interfaces sharing a physical USB parent are grouped together. Other transports use their HID service identity, with a shared Bluetooth HID parent when available. A keyboard exposing independent services may appear more than once; choose its keyboard interface. A USB receiver that presents several keyboards as one HID device cannot distinguish those keyboards.

macOS can reject exclusive access to protected built-in keyboards, and another application may already own the device. The tester reports these failures. Support depends on the OS and device; a normal build cannot bypass an OS-level restriction.

## Controls

During capture, hold **Control + Option** (either left or right), then press:

| Key | Action |
| --- | --- |
| Q | Quit |
| R | Reset all current results |
| P | Release capture and return to device selection |
| H | Toggle help |
| ↑ / ↓ | Scroll the key map and unknown-key list |

The **Devices**, **Reset**, **Help**, and **Quit** buttons also work with a mouse. Scroll the content with the mouse wheel. Controls remain available while help is open, and testing continues. Ordinary letters, Escape, and Ctrl+C on the selected keyboard are testable; Ctrl+C from another keyboard provides an additional way to quit.

In device selection: **↑/↓** select, **Enter** starts, **R** refreshes attached keyboards, **?** opens help, and **Q** or **Escape** quits. Resuming the same device preserves results; selecting a different device clears them.

Normal exit and handled termination signals release capture. Disconnects and input-queue overflow also end capture and show a recovery message. No preferences need restoration because none are modified. All results live in RAM; there are no report files, and quitting discards the results.

## Reference layouts

Percentages describe keyboard sizes, not an exact standard inventory. These are explicit US ANSI reference presets, with rows listed in `internal/keyboard/presets.go`. Rows show key inventory and approximate grouping; tile widths do not reproduce physical keycap sizes. Narrow terminals wrap rows and put the unknown-key panel below the map; wide terminals show it at the side. A 64-column × 18-row terminal is the minimum, and larger sizes reduce scrolling.

| Size | Expected keys | Reference inventory |
| --- | ---: | --- |
| 40% | 44 | Staggered ANSI-inspired compact layout: letters, basic punctuation, eight bottom-row keys; no number row, Caps Lock, or arrows |
| 60% | 61 | Compact ANSI with Escape at top left and no dedicated grave key |
| 65% | 68 | 60% plus Delete, Page Up/Down, and four arrows |
| 75% | 84 | 65% plus a function row, dedicated grave, Insert, Home, and End |
| 80% | 87 | ANSI tenkeyless with function row, navigation cluster, Print Screen, Scroll Lock, and Pause |
| 96% | 100 | Compact full-size inventory: 100% without Insert, Print Screen, Scroll Lock, and Pause |
| 100% | 104 | Full ANSI keyboard with 17 distinct keypad keys |

The GUI/Windows keys are labelled **LCmd/RCmd**, and Alt keys are **LOpt/ROpt**. Known keys outside the preset still appear in the unknown/outside-preset list. The list preserves each HID usage page and usage ID, keeping otherwise similar codes distinct.

The **Last key** line above the map shows the most recently recorded press and its HID code. For example, keypad `*` reports `0007:0055`; Shift+8 is a different physical key combination. If a key does not mark its tile, check whether this line or the press counter changes and whether the key appears in the unknown list. A missing event cannot be corrected by changing the display label.

There is no hold-duration or repetition test. Firmware-only Fn/layer keys, hardware power actions, and other inputs that never reach the HID interface cannot be recorded or suppressed by this application. If a layer changes what a physical key reports, the tester can only identify the emitted usage. Media keys are recorded when exposed by a captured consumer-control interface.

## Development

```sh
make test    # tests with the race detector
make check   # go vet and race tests
```

Tests cover reference inventories, physical-key identity, distinct modifiers, press-edge deduplication, unknown keys, reset, selected-device transitions, shortcuts, permission failures, disconnects, mouse release controls, and responsive rendering. An AST check enforces the hard 30-line limit on Go functions, including tests.

The native HID bridge runs on one locked OS thread. Capture messages carry a generation ID so queued events from an earlier capture cannot enter a later session. Unit tests use a fake capture driver and never seize a real keyboard. Hardware capture and shortcut suppression require a manual test with permission granted.

Implementation is split into `internal/keyboard` for presets, mappings, and results; `internal/capture` for native device access; `internal/ui` for the interface; and `internal/cli` for arguments and application startup.
