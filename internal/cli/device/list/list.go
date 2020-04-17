package list

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/bensallen/rbd/pkg/krbd"
)

const usageHeader = `list - List connected devices

`

// Usage of the list subcommand
func Usage() {
	fmt.Fprintf(os.Stderr, usageHeader)
}

// Run the list subcommand of device
func Run(args []string, verbose bool, noop bool) error {
	devices, err := krbd.Devices()
	if err != nil {
		return err
	}
	if len(devices) != 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "id\tpool\tnamespace\timage\tsnap\tdevice")
		for _, device := range devices {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n", device.ID, device.Pool, device.Namespace, device.Image, device.Snapshot, "/dev/rbd"+strconv.FormatInt(device.ID, 10))
		}
		w.Flush()
	}
	return nil
}
