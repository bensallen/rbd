package unmap

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"testing/iotest"

	"github.com/bensallen/rbd/pkg/krbd"
	flag "github.com/spf13/pflag"
)

const usageHeader = `unmap - Unmap RBD Image

Usage:
  unmap

Flags:
`

var (
	flags = flag.NewFlagSet("unmap", flag.ContinueOnError)
	devid = flags.IntP("devid", "d", -1, "RBD Device ID")
	force = flags.BoolP("force", "f", false, "Optional force argument will wait for running requests and then unmap the image")
)

// Usage of the unmap subcommand
func Usage() {
	fmt.Fprintf(os.Stderr, usageHeader)
	fmt.Fprintf(os.Stderr, flags.FlagUsagesWrapped(0)+"\n")
}

//Usage of rbd with error message
func usageErr(err error) {
	Usage()
	fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
}

// Run the unmap subcommand
func Run(args []string, verbose bool, noop bool) error {
	flags.ParseErrorsWhitelist.UnknownFlags = true
	if err := flags.Parse(args); err != nil {
		Usage()
		fmt.Printf("Error: %v\n\n", err)
		os.Exit(2)
	}

	if *devid == -1 {
		usageErr(errors.New("Device ID not specified"))
		os.Exit(2)
	}

	wc, err := krbd.RBDBusRemoveWriter()
	defer wc.Close()
	w := wc.(io.Writer)
	if err != nil {
		return err
	}

	if verbose {
		w = iotest.NewWriteLogger("unmap", w)
	}

	i := krbd.Image{
		DevID: *devid,
		Options: &krbd.Options{
			Force: *force,
		},
	}

	if noop {
		log.Printf("%s", i)
	} else {
		return i.Unmap(w)
	}
	return nil
}
