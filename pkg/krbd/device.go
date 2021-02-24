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
	Namespace string `krbd:"pool_ns,optional"`
	Image     string `krbd:"name"`
	Snapshot  string `krbd:"current_snap,optional"`
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
		tag := getDeviceTag(t.Field(i))
		if tag.name != "" {
			r, err := os.Open(path + "/" + tag.name)
			if err != nil {
				if tag.optional && errors.Is(err, os.ErrNotExist) {
					continue
				}
				return err
			}
			defer r.Close()
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

// DevPath returns the string form of the Device expected device path, eg. /dev/rbd0
// Does not validate that the device actually exists.
func (d *Device) DevPath() string {
	return "/dev/rbd" + strconv.FormatInt(d.ID, 10)
}

// tag parsing

type deviceTag struct {
	name     string
	optional bool
}

func getDeviceTag(f reflect.StructField) (d deviceTag) {
	s := f.Tag.Get("krbd")
	if s == "" {
		return
	}
	ss := strings.Split(s, ",")
	d.name = ss[0]
	for i := 1; i < len(ss); i++ {
		switch ss[i] {
		case "optional":
			d.optional = true
		}
	}
	return
}
