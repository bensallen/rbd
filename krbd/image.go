package krbd

import (
	"fmt"
	"strings"
)

// Image is a Ceph RBD image.
type Image struct {
	DevID    int // Unmap only
	Monitors []string
	Pool     string
	Image    string
	Snapshot string
	Options  *Options
}

// String mashalls the Image attributes into the string format expected by the krbd add interface, eg:
// "${mons} name=${user},secret=${key} ${pool} ${image} ${snap}"
func (i Image) String() string {
	if i.Snapshot == "" {
		i.Snapshot = "-"
	}
	return fmt.Sprintf("%s %s %s %s %s", strings.Join(i.Monitors, ","), i.Options, i.Pool, i.Image, i.Snapshot)
}
