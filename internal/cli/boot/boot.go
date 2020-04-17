package boot

import (
	"fmt"
	"log"
	"os"
	"testing/iotest"

	"github.com/bensallen/rbd/pkg/cmdline"
	"github.com/bensallen/rbd/pkg/krbd"
	flag "github.com/spf13/pflag"
)

const usageHeader = `boot - Boot via RBD image

Usage:
  boot

Flags:
`

var (
	flags    = flag.NewFlagSet("boot", flag.ContinueOnError)
	procPath = flags.StringP("cmdline", "c", "/proc/cmdline", "Path to kernel cmdline (default: /proc/cmdline)")
)

// Usage of the boot subcommand
func Usage() {
	fmt.Fprintf(os.Stderr, usageHeader)
	fmt.Fprintf(os.Stderr, flags.FlagUsagesWrapped(0)+"\n")
}

// Run the boot subcommand
func Run(args []string, verbose bool, noop bool) error {
	flags.ParseErrorsWhitelist.UnknownFlags = true
	if err := flags.Parse(args); err != nil {
		Usage()
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		os.Exit(2)
	}
	procCmdline, err := cmdline.Read(*procPath)
	if err != nil {
		return err
	}

	mounts := cmdline.Parse(string(procCmdline))

	w, err := krbd.RBDBusAddWriter()
	if err != nil {
		return err
	}

	if verbose {
		w = iotest.NewWriteLogger("Boot:", w)
	}

	for name, mount := range mounts {
		log.Printf("Boot: mapping image %s from cmdline", name)

		if noop {
			log.Printf("%s", mount.Image)
		} else {
			err = mount.Image.Map(w)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
