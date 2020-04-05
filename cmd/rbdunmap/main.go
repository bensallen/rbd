package main

import (
	"log"
	"testing/iotest"

	"github.com/bensallen/rbdmap/krbd"
	flag "github.com/spf13/pflag"
)

var (
	devid   = flag.Int("devid", 0, "RBD Device Index (default 0)")
	force   = flag.Bool("force", false, "Optional force argument will wait for running requests and then unmap the image")
	dryRun  = flag.BoolP("dry-run", "n", false, "dry run (don't actually unmap device)")
	verbose = flag.BoolP("verbose", "v", false, "Enable additional output")
)

func main() {
	flag.Parse()

	w, err := krbd.RBDBusRemoveWriter()
	if err != nil {
		log.Fatalf("%#v", err)
	}

	if *verbose {
		w = iotest.NewWriteLogger("unmap", w)
	}

	i := krbd.Image{
		DevID: *devid,
		Options: &krbd.Options{
			Force: *force,
		},
	}

	if *dryRun {
		log.Printf("%s", i)
	} else {
		err = i.Unmap(w)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}
