package main

import (
	"log"
	"testing/iotest"

	"github.com/bensallen/rbdmap/krbd"
	flag "github.com/spf13/pflag"
)

var (
	devid   = flag.Int("devid", 0, "RBD Device Index (default 0")
	force   = flag.Bool("force", false, "Optional force argument will wait for running requests and then unmap the image")
	dryRun  = flag.Bool("dry-run", false, "dry run (don't actually unmap image)")
	verbose = flag.Bool("verbose", false, "Enable additional output")
	help    = flag.Bool("help", false, "Display usage")
)

func main() {
	flag.Parse()

	//if *device == "" {
	//	fmt.Printf("Error: --device must be specified\n")
	//	flag.PrintDefaults()
	//	os.Exit(1)
	//}

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
