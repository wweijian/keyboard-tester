package cli

import (
	"flag"
	"fmt"
	"io"
)

type options struct {
	keyboard string
	list     bool
}

func parse(args []string, output io.Writer) (options, error) {
	var opts options
	flags := flag.NewFlagSet("keyboard-tester", flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&opts.keyboard, "keyboard", "", "keyboard size: 40, 60, 65, 75, 80, 96, 100 (optional % suffix)")
	flags.BoolVar(&opts.list, "list", false, "list attached keyboards without starting capture")
	flags.Usage = func() { help(output) }
	if err := flags.Parse(args); err != nil {
		return opts, err
	}
	if flags.NArg() != 0 {
		return opts, fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}
	return opts, nil
}

func help(output io.Writer) {
	fmt.Fprint(output, `Keyboard Lab — macOS physical keyboard tester

Usage:
  sudo keyboard-tester --keyboard 65
  sudo keyboard-tester --keyboard 100%
  keyboard-tester --list

Options:
  --keyboard SIZE   40%, 60%, 65%, 75%, 80%, 96%, 100% (the % is optional)
  --list            List connected keyboards without capturing input
  -h, --help        Show help

Select a physical keyboard in the UI, then press Enter to start.
Capture requires administrator privileges (sudo), Input Monitoring permission,
and exclusive device access. Help and --list do not require sudo.
During capture: Ctrl+Option+Q quit, +R reset, +P devices, +H help.
Mouse buttons also provide controls. Results are kept only in RAM.
`)
}
