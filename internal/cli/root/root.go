package root

import (
	"errors"
	"fmt"
	"os"

	"github.com/bensallen/rbdmap/internal/cli/boot"
	"github.com/bensallen/rbdmap/internal/cli/device"
	"github.com/bensallen/rbdmap/internal/cli/device/list"
	"github.com/bensallen/rbdmap/internal/cli/rbdmap"
	"github.com/bensallen/rbdmap/internal/cli/unmap"
	flag "github.com/spf13/pflag"
)

// Version is the CLI version or release number. To be overriden at build time.
var Version = "unknown"

var (
	rootFlags = flag.NewFlagSet("root", flag.ContinueOnError)
	help      = rootFlags.BoolP("help", "h", false, "Diplay help.")
	version   = rootFlags.BoolP("version", "V", false, "Displays the program version string.")
	noop      = rootFlags.BoolP("noop", "n", false, "No-op (don't actually perform action).")
	verbose   = rootFlags.BoolP("verbose", "v", false, "Enable additional output.")
)

const usageHeader = `rbd - Ceph RBD CLI

Usage:
  rbd [map|unmap|device|boot]

Subcommands:
  map      Map RBD Image
  unmap    Unmap RBD Image
  boot     Boot via RBD Image
  device   Manage RBD Devices

Flags:
`

//Usage of rbd command printed to stderr
func usage() {
	fmt.Fprintf(os.Stderr, usageHeader)
	fmt.Fprintf(os.Stderr, rootFlags.FlagUsagesWrapped(0)+"\n")
}

//Usage of rbd with error message
func usageErr(err error) {
	usage()
	fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
}

// Run rbd command
func Run(args []string) error {
	rootFlags.ParseErrorsWhitelist.UnknownFlags = true
	if err := rootFlags.Parse(args); err != nil {
		usageErr(err)
		os.Exit(2)
	}

	if *version {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	// Pick the correct usage to call based on subcommand
	if *help {
		if rootFlags.NArg() == 0 {
			usage()
		} else {
			switch rootFlags.Arg(0) {
			case "map":
				rbdmap.Usage()
			case "unmap":
				unmap.Usage()
			case "device":
				// device has its own subcommands
				switch rootFlags.Arg(1) {
				case "list":
					list.Usage()
				case "map":
					rbdmap.Usage()
				case "unmap":
					unmap.Usage()
				default:
					device.Usage()
				}
			case "boot":
				boot.Usage()
			}
		}
		os.Exit(0)
	}

	if rootFlags.NArg() < 1 {
		usageErr(errors.New("missing subcommand"))
		os.Exit(2)
	}

	// Run subcommands
	switch rootFlags.Arg(0) {
	case "map":
		return rbdmap.Run(args, *verbose, *noop)
	case "unmap":
		return unmap.Run(args, *verbose, *noop)
	case "device":
		return device.Run(args, *verbose, *noop)
	case "boot":
		return boot.Run(args, *verbose, *noop)
	case "help":
		usage()
	default:
		usageErr(fmt.Errorf("unrecognized subcommand: %s", rootFlags.Arg(0)))
		os.Exit(2)
	}
	return nil
}
