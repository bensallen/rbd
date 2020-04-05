package krbd

import (
	"fmt"
	"reflect"
	"strings"
)

// Options is per client instance and per mapping (block device) rbd device map options.
// krbd tag is the string of the option passed via sysfs.
// Reference: https://docs.ceph.com/docs/master/man/8/rbd/#kernel-rbd-krbd-options
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

// String marshalls Options via the krbd struct tags into comma seperated format
// that matches the format expected via the krbd add interface.
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
		// Bool types don't include their value just the tag.
		if v.Field(i).Kind() == reflect.Bool {
			output = append(output, tag)
		} else {
			output = append(output, fmt.Sprintf("%s=%v", tag, v.Field(i)))
		}
	}
	return strings.Join(output, ",")
}
