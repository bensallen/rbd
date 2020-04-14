package main

import (
	"log"
	"testing/iotest"

	"github.com/bensallen/rbdmap/cmdline"
	"github.com/bensallen/rbdmap/krbd"
	flag "github.com/spf13/pflag"
)

var (
	procPath = flag.StringP("cmdline", "c", "/proc/cmdline", "Path to kernel cmdline (default: /proc/cmdline)")
	dryRun   = flag.BoolP("dry-run", "n", false, "dry run (don't actually map image)")
	verbose  = flag.BoolP("verbose", "v", false, "Enable additional output")
)

func main() {
	flag.Parse()

	procCmdline, err := cmdline.Read(*procPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	mounts := cmdline.Parse(string(procCmdline))

	w, err := krbd.RBDBusAddWriter()
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *verbose {
		w = iotest.NewWriteLogger("boot", w)
	}

	for name, mount := range mounts {
		log.Printf("Mapping image %s from cmdline", name)

		if !*dryRun {
			err = mount.Image.Map(w)
			if err != nil {
				log.Fatalf("%v", err)
			}
		}
	}
}
