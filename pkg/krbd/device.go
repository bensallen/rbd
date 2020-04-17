package krbd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Device is a RBD device discovered via /sys/bus/rbd/device and attributes read
// from /sys/device/rbd/<id>. Krbd struct tags are the sysfs file names that
// are used, eg. for Pool - /sys/device/rbd/0/pool
type Device struct {
	ID        int64
	Pool      string `krbd:"pool"`
	Namespace string `krbd:"pool_ns"`
	Image     string `krbd:"name"`
	Snapshot  string `krbd:"current_snap"`
}

// Devices iterates over /sys/bus/rbd/device/ to find all mapped RBD devices populating
// attributes.
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

		id, err := strconv.ParseInt(p.Name(), 10, 0)
		if err != nil {
			return nil, err
		}

		d := Device{ID: id}
		if err := d.readDeviceAttrs(realPath); err != nil {
			return nil, err
		}
		devices[i] = d
	}

	// Sort devices based on ID
	sort.SliceStable(devices, func(i, j int) bool {
		return devices[i].ID < devices[j].ID
	})

	return devices, nil
}

func (d *Device) readDeviceAttrs(path string) error {
	t := reflect.TypeOf(*d)
	v := reflect.ValueOf(d).Elem()

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
			if len(value) != 0 {
				v.Field(i).SetString(strings.TrimSpace(string(value)))
			}
		}
	}
	return nil
}

// Find looks for an existing mapped device based on the set attributes
// of (d *Device), doesn't match on ID. If a match is found from sysfs
// (d *Device) remaining attributes are updated from sysfs.
func (d *Device) Find() error {
	devices, err := Devices()
	if err != nil {
		return err
	}
	return d.find(devices)
}

func (d *Device) find(devices []Device) error {
	empty := Device{}
	if *d == empty {
		return errors.New("Device has no attributes set")
	}
	for _, device := range devices {
		if d.Image != "" && device.Image != d.Image {
			continue
		}
		if d.Namespace != "" && device.Namespace != d.Namespace {
			continue
		}
		if d.Snapshot != "" && device.Snapshot != d.Snapshot {
			continue
		}
		if d.Pool != "" && device.Pool != d.Pool {
			continue
		}
		*d = device
		return nil
	}
	return errors.New("No match found")
}
