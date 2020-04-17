package krbd

import (
	"errors"
	"fmt"
	"io"
)

// Map the RBD image via the krbd interface. An open io.Writer is required
// typically to /sys/bus/rbd/add or /sys/bus/rbd/add_single_major
func (i *Image) Map(w io.Writer) error {
	if len(i.Monitors) == 0 {
		return errors.New("No monitors defined")
	}
	if i.Pool == "" {
		return errors.New("No pool defined")
	}
	if i.Image == "" {
		return errors.New("No image defined")
	}

	out := i.String()
	n, err := w.Write([]byte(out))

	if n != len(out) {
		return fmt.Errorf("Incomplete write, wrote %d, expected to write %d", n, len(out))
	}
	return err
}
