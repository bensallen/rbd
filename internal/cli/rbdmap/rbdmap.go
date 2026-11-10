package rbdmap

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing/iotest"

	"github.com/bensallen/rbd/pkg/krbd"
	flag "github.com/spf13/pflag"
)

const usageHeader = `map - Map RBD Image

Usage:
  map

Flags:
`

var (
	flags     = flag.NewFlagSet("map", flag.ContinueOnError)
	monAddrs  = flags.StringSliceP("monitor", "m", []string{}, "Connect to one or more monitor addresses (192.168.0.1[:6789]). Multiple address are specified comma separated.")
	pool      = flags.StringP("pool", "p", "", "Interact with the given pool.")
	image     = flags.StringP("image", "i", "", "Image to map")
	namespace = flags.String("namespace", "", "Use a pre-defined image namespace within a pool")
	snap      = flags.String("snap", "", "Specifies a snapshot name")
	id        = flags.String("id", "", "Specifies the username (without the 'client.' prefix)")
	secret    = flags.String("secret", "", "Specifies the user authentication secret")
	readOnly  = flags.Bool("read-only", false, "Map the image read-only")
	msMode    = flags.String("ms-mode", "", "Messenger mode: legacy, crc, secure, prefer-crc, or prefer-secure (msgr2 requires crc or secure)")
)

// Usage of the map subcommand
func Usage() {
	fmt.Fprintf(os.Stderr, usageHeader)
	fmt.Fprint(os.Stderr, flags.FlagUsagesWrapped(0)+"\n")
}

// Run the map subcommand
func Run(args []string, verbose bool, noop bool) error {
	flags.ParseErrorsWhitelist.UnknownFlags = true
	if err := flags.Parse(args); err != nil {
		Usage()
		fmt.Printf("Error: %v\n\n", err)
		os.Exit(2)
	}

	if len(*monAddrs) == 0 || *pool == "" || *image == "" || *id == "" || *secret == "" {
		Usage()
		fmt.Print("Error: --monitor, --pool, --image, --id, and --secret must be specified\n\n")
		os.Exit(2)
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
			MsMode:    *msMode,
		},
	}

	if noop {
		log.Printf("%s", i)
		return nil
	}

	wc, err := krbd.RBDBusAddWriter()
	if err != nil {
		return err
	}
	defer wc.Close()
	w := io.Writer(wc)

	if verbose {
		w = iotest.NewWriteLogger("map", w)
	}

	return i.Map(w)
}
