package device

import (
	"errors"
	"fmt"
	"os"

	"github.com/bensallen/rbdmap/internal/cli/rbdmap"
	"github.com/bensallen/rbdmap/internal/cli/unmap"
	flag "github.com/spf13/pflag"
)

const usageHeader = `device - Manage RBD Devices

Usage:
  device [list|map|unmap]

Subcommands:
  list     List connected devices
  map      Map RBD Image
  unmap    Unmap RBD Image

`

var (
	flags = flag.NewFlagSet("unmap", flag.ContinueOnError)
)

// Usage of the unmap subcommand
func Usage() {
	fmt.Fprintf(os.Stderr, usageHeader)
	//fmt.Fprintf(os.Stderr, flags.FlagUsagesWrapped(0)+"\n")
}

//Usage of rbd with error message
func usageErr(err error) {
	Usage()
	fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
}

// Run the device subcommand
func Run(args []string, verbose bool, noop bool) error {
	flags.ParseErrorsWhitelist.UnknownFlags = true
	if err := flags.Parse(args); err != nil {
		Usage()
		fmt.Printf("Error: %v\n\n", err)
		os.Exit(2)
	}

	if flags.NArg() < 2 {
		usageErr(errors.New("missing subcommand"))
		os.Exit(2)
	}

	switch flags.Arg(1) {
	case "list":
		// return list.Run(args, verbose, noop)
	case "map":
		return rbdmap.Run(args, verbose, noop)
	case "unmap":
		return unmap.Run(args, verbose, noop)
	case "help":
		Usage()
	default:
		usageErr(fmt.Errorf("unrecognized subcommand: %s", flags.Arg(0)))
		os.Exit(2)
	}
	return nil
}
