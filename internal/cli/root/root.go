package root

import (
	"fmt"
	"os"

	"github.com/bensallen/rbdmap/internal/cli/boot"
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

// Usage help text
const usageHeader = `rbd - Ceph RBD CLI

Usage:
  rbd [map|unmap|device|boot]

Subcommands:
  map      Map RBD Image
  unmap    Unmap RBD Image
  device   RBD device
  boot     Boot via RBD image

Flags:
`

//Usage of rbd command
func Usage() {
	fmt.Fprintf(os.Stderr, usageHeader)
	fmt.Fprintf(os.Stderr, rootFlags.FlagUsagesWrapped(0)+"\n")
}

// Run rbd command
func Run(args []string) error {
	rootFlags.ParseErrorsWhitelist.UnknownFlags = true
	if err := rootFlags.Parse(args); err != nil {
		Usage()
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		os.Exit(2)
	}

	if *version {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	if *help {
		if rootFlags.NArg() == 0 {
			Usage()
		} else {
			switch rootFlags.Arg(0) {
			case "map":
				rbdmap.Usage()
			case "unmap":
				unmap.Usage()
			case "device":

			case "boot":
				boot.Usage()
			}
		}
		os.Exit(0)
	}

	if rootFlags.NArg() < 1 {
		return fmt.Errorf("Missing subcommand")
	}

	if rootFlags.NArg() > 1 {
		return fmt.Errorf("Multiple subcommands specified")
	}

	switch rootFlags.Arg(0) {
	case "map":
		return rbdmap.Run(args, *verbose, *noop)
	case "unmap":
		return unmap.Run(args, *verbose, *noop)
	case "device":

	case "boot":
		return boot.Run(args, *verbose, *noop)
	default:
		Usage()
		fmt.Fprintf(os.Stderr, "Error: unrecognized subcommand: %s\n", rootFlags.Arg(0))
		os.Exit(2)
	}
	return nil
}
