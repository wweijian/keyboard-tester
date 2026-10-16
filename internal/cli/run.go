package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"

	"keyboard-tester/internal/capture"
	"keyboard-tester/internal/keyboard"
	"keyboard-tester/internal/ui"
)

func Run(args []string, output, errorOutput io.Writer) error {
	if len(args) == 0 {
		help(output)
		return nil
	}
	opts, err := parse(args, output)
	if errors.Is(err, flag.ErrHelp) {
		return nil
	}
	if err != nil {
		return err
	}
	return execute(opts, output, errorOutput)
}

func execute(opts options, output, errorOutput io.Writer) error {
	if opts.list {
		return listDevices(output)
	}
	layout, err := keyboard.Lookup(opts.keyboard)
	if err != nil {
		return err
	}
	if !term.IsTerminal(os.Stdin.Fd()) {
		return errors.New("interactive testing requires a terminal; use --list to list devices")
	}
	if err := requireAdmin(layout.Size, os.Geteuid()); err != nil {
		return err
	}
	driver, err := capture.New()
	if err != nil {
		return err
	}
	defer driver.Close()
	return interactive(driver, layout, output, errorOutput)
}

func requireAdmin(size string, uid int) error {
	if uid == 0 {
		return nil
	}
	return fmt.Errorf("exclusive keyboard capture requires administrator privileges. Run: sudo ./bin/keyboard-tester --keyboard %s", size)
}

func interactive(driver capture.Driver, layout keyboard.Layout, output, errorOutput io.Writer) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTSTP)
	defer stop()
	program := tea.NewProgram(ui.New(driver, layout), tea.WithOutput(output), tea.WithContext(ctx))
	_, err := program.Run()
	if ctx.Err() != nil {
		fmt.Fprintln(errorOutput, "Keyboard capture released.")
		return nil
	}
	return err
}

func listDevices(output io.Writer) error {
	driver, err := capture.New()
	if err != nil {
		return err
	}
	defer driver.Close()
	devices, err := driver.Devices()
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		fmt.Fprintln(output, "No keyboards found.")
	}
	for _, device := range devices {
		fmt.Fprintf(output, "%X  %s  [%s]  %04X:%04X\n", device.ID, device.Name, device.Transport, device.Vendor, device.Product)
	}
	return nil
}
