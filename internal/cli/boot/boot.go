package boot

import (
	"fmt"
	"log"
	"os"
	"testing/iotest"

	"github.com/bensallen/rbd/pkg/cmdline"
	"github.com/bensallen/rbd/pkg/krbd"
	"github.com/bensallen/rbd/pkg/mount"
	flag "github.com/spf13/pflag"
)

const usageHeader = `boot - Boot via RBD image

Usage:
  boot

Flags:
`

var (
	flags    = flag.NewFlagSet("boot", flag.ContinueOnError)
	mkdir    = flags.BoolP("mkdir", "m", false, "Create the destination mount path if it doesn't exist")
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

	for name, mnt := range mounts {
		log.Printf("Boot: mapping image %s from cmdline", name)

		if noop {
			log.Printf("%s", mnt.Image)
		} else {

			// Map the RBD device
			if err := mnt.Image.Map(w); err != nil {
				return err
			}

			// Find the RBD device that was just mapped
			dev := krbd.Device{Image: mnt.Image.Image, Pool: mnt.Image.Pool, Namespace: mnt.Image.Options.Namespace, Snapshot: mnt.Image.Snapshot}
			if err := dev.Find(); err != nil {
				return err
			}

			if verbose {
				log.Printf("Boot: device found %#v\n", dev)
			}

			if *mkdir {
				if err := os.MkdirAll(mnt.Path, 0755); err != nil {
					return err
				}
			}

			// Attempt to mount the device
			return mount.Mount(dev.DevPath(), mnt.Path, mnt.FsType, mnt.MountOpts)
		}
	}
	return nil
}
