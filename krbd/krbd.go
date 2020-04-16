package krbd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
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

type Device struct {
	ID        int
	Pool      string `krbd:"pool"`
	Namespace string `krbd:"pool_ns"`
	Image     string `krbd:"name"`
	Snapshot  string `krbd:"current_snap"`
}

func Devices() ([]Device, error) {
	path := sysBusPath + "/devices"
	return devices(path)
}

func devices(path string) ([]Device, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	list, err := f.Readdir(-1)
	f.Close()

	if err != nil {
		return nil, err
	}

	devices := make([]Device, len(list))

	for i, p := range list {
		realPath, err := filepath.EvalSymlinks(path + "/" + p.Name())
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(p.Name())
		if err != nil {
			return nil, err
		}

		d := Device{ID: id}
		if err := d.readDeviceAttrs(realPath); err != nil {
			return nil, err
		}
		devices[i] = d
	}
	return devices, nil
}

func (d *Device) readDeviceAttrs(path string) error {
	t := reflect.TypeOf(d)
	v := reflect.ValueOf(d)

	// Iterate over all available struct fields
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("krbd")
		if tag != "" {
			r, err := os.Open(path + "/" + tag)
			defer r.Close()
			if err != nil {
				return err
			}
			value, err := ioutil.ReadAll(r)

			if err != nil {
				return err
			}
			v.Field(i).SetString(string(value))
		}
	}
	return nil
}
