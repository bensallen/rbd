package krbd

import (
	"fmt"
	"io"
	"os"
)

// Reference: https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-bus-rbd

const sysBusPath = "/sys/bus/rbd"

// RBDBusAddWriter returns an io.Writer with the appropriate sysfs rbd/add opened.
func RBDBusAddWriter() (io.Writer, error) {
	rbdBusAddSingleMajor := sysBusPath + "/add_single_major"
	rbdBusAdd := sysBusPath + "/add"

	if _, err := os.Stat(rbdBusAddSingleMajor); err == nil {
		return os.OpenFile(rbdBusAddSingleMajor, os.O_WRONLY, 0644)
	} else if _, err := os.Stat(rbdBusAdd); err == nil {
		return os.OpenFile(rbdBusAdd, os.O_WRONLY, 0644)
	}
	return nil, fmt.Errorf("could not find %s or %s", rbdBusAddSingleMajor, rbdBusAdd)
}

// RBDBusRemoveWriter returns an io.Writer with the appropriate sysfs rbd/remove opened.
func RBDBusRemoveWriter() (io.Writer, error) {
	rbdBusAddSingleMajor := sysBusPath + "/remove_single_major"
	rbdBusAdd := sysBusPath + "/remove"

	if _, err := os.Stat(rbdBusAddSingleMajor); err == nil {
		return os.OpenFile(rbdBusAddSingleMajor, os.O_WRONLY, 0644)

	} else if _, err := os.Stat(rbdBusAdd); err == nil {
		return os.OpenFile(rbdBusAdd, os.O_WRONLY, 0644)
	}
	return nil, fmt.Errorf("could not find %s or %s", rbdBusAddSingleMajor, rbdBusAdd)
}
