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
	flags      = flag.NewFlagSet("boot", flag.ContinueOnError)
	mkdir      = flags.BoolP("mkdir", "m", false, "Create the destination mount path if it doesn't exist")
	switchRoot = flags.BoolP("switch-root", "s", false, "Attempt to switch_root to root filesystem and execute /sbin/init")
	procPath   = flags.StringP("cmdline", "c", "/proc/cmdline", "Path to kernel cmdline (default: /proc/cmdline)")
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

			if mnt.Path == "" {
				return fmt.Errorf("device path not set")
			}

			if mnt.Path[0] != '/' {
				return fmt.Errorf("device path not an absolute path")
			}

			if *mkdir {
				if err := os.MkdirAll(mnt.Path, 0755); err != nil {
					return err
				}
			}

			// Attempt to mount the device
			if err := mount.Mount(dev.DevPath(), mnt.Path, mnt.FsType, mnt.MountOpts); err != nil {
				return err
			}
		}
	}

	if root, ok := mounts["root"]; ok {
		if root.Overlay {
			if verbose {
				log.Printf("Boot: attempting to mount overlay over %s\n", root.Path)
			}
			if !noop {
				if err := mount.Overlay(root.Path); err != nil {
					return err
				}
			}
		}
		if *switchRoot {
			if verbose {
				log.Printf("Boot: attempting to switch root to %s with /sbin/init\n", root.Path)
			}
			if !noop {
				return mount.SwitchRoot(root.Path, "/sbin/init")
			}
		}
	}
	return nil
}
