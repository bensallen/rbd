package krbd

import (
	"fmt"
	"io"
	"os"
)

const sysBusPath = "/sys/bus/rbd"

// RBDBusAdd returns an io.Writer with the appropriate sysfs rbd/add opened.
func RBDBusAdd() (io.Writer, error) {
	rbdBusAddSingleMajor := sysBusPath + "/add_single_major"
	rbdBusAdd := sysBusPath + "/add"

	if _, err := os.Stat(rbdBusAddSingleMajor); err == nil {
		return os.OpenFile(rbdBusAddSingleMajor, os.O_WRONLY, 0644)

	} else if _, err := os.Stat(rbdBusAdd); err == nil {
		return os.OpenFile(rbdBusAdd, os.O_WRONLY, 0644)

	}
	return nil, fmt.Errorf("Both %s and %s do not exist", rbdBusAddSingleMajor, rbdBusAdd)
}
