package krbd

import (
	"io"
	"strconv"
)

// Unmap a RBD device via the krbd interface. DevID must be defined.
// An open io.Writer is required typically to /sys/bus/rbd/remove
// or /sys/bus/rbd/remove_single_major.
func (i *Image) Unmap(w io.Writer) error {

	cmd := strconv.Itoa(i.DevID)
	if i.Options != nil && i.Options.Force {
		cmd = cmd + " force"
	}
	_, err := w.Write([]byte(cmd))
	return err
}
