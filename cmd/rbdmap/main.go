package main

import (
	"fmt"
	"log"
	"os"
	"testing/iotest"

	"github.com/bensallen/rbdmap/krbd"
	flag "github.com/spf13/pflag"
)

var (
	monAddrs  = flag.StringSliceP("monitor", "m", []string{}, "Connect to one or more monitor addresses (192.168.0.1[:6789]). Multiple address are specified comma separated.")
	pool      = flag.StringP("pool", "p", "", "Interact with the given pool.")
	image     = flag.StringP("image", "i", "", "Image to map")
	namespace = flag.String("namespace", "", "Use a pre-defined image namespace within a pool")
	snap      = flag.String("snap", "", "Specifies a snapshot name")
	id        = flag.String("id", "", "Specifies the username (without the 'client.' prefix)")
	secret    = flag.String("secret", "", "Specifies the user authentication secret")
	readOnly  = flag.Bool("read-only", false, "Map the image read-only")
	dryRun    = flag.BoolP("dry-run", "n", false, "dry run (don't actually map image)")
	verbose   = flag.BoolP("verbose", "v", false, "Enable additional output")
)

func main() {
	flag.Parse()

	if len(*monAddrs) == 0 || *pool == "" || *image == "" || *id == "" || *secret == "" {
		fmt.Printf("Error: --monitor, --pool, --image, --id, and --secret must be specified\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	w, err := krbd.RBDBusAddWriter()
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *verbose {
		w = iotest.NewWriteLogger("map", w)
	}

	i := krbd.Image{
		Monitors: *monAddrs,
		Pool:     *pool,
		Image:    *image,
		Snapshot: *snap,
		Options: &krbd.Options{
			ReadOnly:  *readOnly,
			Name:      *id,
			Secret:    *secret,
			Namespace: *namespace,
		},
	}

	if *dryRun {
		log.Printf("%s", i)
	} else {
		err = i.Map(w)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}
