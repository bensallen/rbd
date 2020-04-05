package krbd

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Image is a Ceph RBD image.
type Image struct {
	DevID    int // Unmap only
	Monitors []string
	Pool     string
	Image    string
	Options  *Options
	Snapshot string
}

func (i Image) String() string {
	if i.Snapshot == "" {
		i.Snapshot = "-"
	}
	// "${mons} name=${user},secret=${key} ${pool} ${image} ${snap}"
	return fmt.Sprintf("%s %s %s %s %s", strings.Join(i.Monitors, ","), i.Options, i.Pool, i.Image, i.Snapshot)
}

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

	_, err := w.Write([]byte(i.String()))
	return err
}

// Unmap a RBD device via the krbd interface. DevID must be defined.
// An open io.Writer is required typically to /sys/bus/rbd/remove
// or /sys/bus/rbd/remove_single_major.
func (i *Image) Unmap(w io.Writer) error {
	if i.DevID == 0 {
		return errors.New("DevID not defined")
	}

	cmd := strconv.Itoa(i.DevID)
	if i.Options != nil && i.Options.Force {
		cmd = cmd + " force"
	}
	_, err := w.Write([]byte(cmd))
	return err
}

// Options is per client instance and per mapping (block device) rbd device map options.
// krbd tag is the string of the option passed via sysfs.
type Options struct {
	// Client Options
	Fsid                     string `krbd:"fsid"`
	IP                       string `krbd:"ip"`
	Share                    bool   `krbd:"share"`
	Noshare                  bool   `krbd:"noshare"`
	CRC                      bool   `krbd:"crc"`
	NoCRC                    bool   `krbd:"nocrc"`
	CephxRequireSignatures   bool   `krbd:"cephx_require_signatures"`
	NoCephxRequireSignatures bool   `krbd:"nocephx_require_signatures"`
	TCPNoDelay               bool   `krbd:"tcp_nodelay"`
	NoTCPNoDelay             bool   `krbd:"notcp_nodelay"`
	CephxSignMessages        bool   `krbd:"cephx_sign_messages"`
	NoCephxSignMessages      bool   `krbd:"nocephx_sign_messages"`
	MountTimeout             int    `krbd:"mount_timeout"`
	OSDKeepAlive             int    `krbd:"osdkeepalive"`
	OSDIdleTTL               int    `krbd:"osd_idle_ttl"`

	// RBD Block Options
	Force       bool   `krbd:"force"` // Unmap only
	ReadWrite   bool   `krbd:"rw"`
	ReadOnly    bool   `krbd:"ro"`
	QueueDepth  int    `krbd:"queue_depth"`
	LockOnRead  bool   `krbd:"lock_on_read"`
	Exclusive   bool   `krbd:"exclusive"`
	LockTimeout uint64 `krbd:"lock_timeout"`
	NoTrim      bool   `krbd:"notrim"`
	AbortOnFull bool   `krbd:"abort_on_full"`
	AllocSize   int    `krbd:"alloc_size"`
	Name        string `krbd:"name"`
	Secret      string `krbd:"secret"`
	Namespace   string `krbd:"_pool_ns"`
}

func (o Options) String() string {
	output := []string{}

	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)

	// Iterate over all available struct fields
	for i := 0; i < t.NumField(); i++ {
		// Skip values that are zero values of the struct. Otherwise Options would have
		// to track the upstream default values to always provide all options.
		if v.Field(i).Interface() == reflect.Zero(v.Field(i).Type()).Interface() {
			continue
		}
		tag := t.Field(i).Tag.Get("krbd")
		if v.Field(i).Kind() == reflect.Bool {
			output = append(output, fmt.Sprintf("%s", tag))
		} else {
			output = append(output, fmt.Sprintf("%s=%v", tag, v.Field(i)))
		}
	}
	return strings.Join(output, ",")
}
