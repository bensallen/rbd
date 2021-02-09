package boot

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing/iotest"

	"github.com/bensallen/rbd/pkg/boot"
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
	flags       = flag.NewFlagSet("boot", flag.ContinueOnError)
	mkdir       = flags.BoolP("mkdir", "m", false, "Create the destination mount path if it doesn't exist")
	switchRoot  = flags.StringP("switch-root", "s", "", "Attempt to switch_root to root filesystem and execute provided init path")
	unshareRoot = flags.StringP("unshare", "u", "", "Attempt to execute init in a namespaced context (container) inside the root filesystem")
	procPath    = flags.StringP("cmdline", "c", "/proc/cmdline", "Path to kernel cmdline (default: /proc/cmdline)")
)

const (
	// RootPath is the path that is prepended for all mounts, and is the path that will be switch_root'ed to
	// if a root mounted is mounted.
	RootPath = "/newroot"

	// OverlayPath is used when the root mount indicates that a overlay should be used. The upper and work directories are
	// created under this directory as OverlayPath + "/upper" and OverlayPath + "/work".
	OverlayPath = "/run/overlayfs"

	// OverlayRootPath is the path that is prepended for all mounts when the root mount indicates that a overlay
	// should be mounted over it. The overlay mounts to RootPath and switch_root still switches to RootPath.
	OverlayRootPath = "/run/overlayfs/lower"
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

	wc, err := krbd.RBDBusAddWriter()
	defer wc.Close()
	w := wc.(io.Writer)
	if err != nil {
		return err
	}

	if verbose {
		w = iotest.NewWriteLogger("Boot:", w)
	}

	// Set the prepended mount path
	mntPrefix := RootPath
	if root, ok := mounts["root"]; ok {
		if root.Path != "/" {
			return fmt.Errorf("root path must be set to '/'")
		}
		if root.Overlay {
			mntPrefix = OverlayRootPath
		}
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
				if err := os.MkdirAll(mntPrefix+mnt.Path, 0755); err != nil {
					return err
				}
			}

			// Attempt to mount the device
			if err := mount.Mount(dev.DevPath(), mntPrefix+mnt.Path, mnt.FsType, mnt.MountOpts); err != nil {
				return err
			}
		}
	}

	if root, ok := mounts["root"]; ok {
		if root.Overlay {
			if verbose {
				log.Printf("Boot: attempting to mount root overlay to %s\n", RootPath)
			}
			if !noop {
				if err := mount.Overlay(OverlayRootPath, OverlayPath+"/upper", OverlayPath+"/work", RootPath); err != nil {
					return err
				}
			}
		}
		if *switchRoot != "" {
			if verbose {
				log.Printf("Boot: attempting to switch root to %s with init %s\n", RootPath, *switchRoot)
			}
			if !noop {
				return boot.SwitchRoot(RootPath, *switchRoot)
			}
		}
		if *unshareRoot != "" {
			if verbose {
				log.Printf("Boot: attempting to execute init in a namespaced context %s with init %s\n", RootPath, *unshareRoot)
			}
			if !noop {
				return boot.UnshareRoot(RootPath, *unshareRoot)
			}
		}
	}
	return nil
}
