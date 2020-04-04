package krbd

import (
	"fmt"
	"reflect"
	"strings"
)

// Image is a Ceph RBD image.
type Image struct {
	DevID     int
	Monitors  []string
	Pool      string
	Namespace string
	Image     string
	Options   *Options
	Snapshot  string
}

// Options to be passed when mapping a RBD image.
// krbd tag is the string of the option passed via sysfs.
type Options struct {
	Exclusive   bool   `krbd:"exclusive"`
	LockOnRead  bool   `krbd:"lock_on_read"`
	NoTrim      bool   `krbd:"notrim"`
	ReadOnly    bool   `krbd:"read_only"`
	AllocSize   int    `krbd:"alloc_size"`
	QueueDepth  int    `krbd:"queue_depth"`
	LockTimeout uint64 `krbd:"lock_timeout"`
	Name        string `krbd:"name"`
	Namespace   string `krbd:"_pool_ns"`
	Secret      string `krbd:"secret"`
}

func (o Options) String() string {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)

	output := []string{}

	// Iterate over all available struct fields
	for i := 0; i < t.NumField(); i++ {

		// Skip values that are zero values of the struct. Otherwise Options would have
		// to track the upstream default values and always provide all options.
		if v.Field(i).Interface() == reflect.Zero(v.Field(i).Type()).Interface() {
			continue
		}
		tag := t.Field(i).Tag.Get("krbd")
		output = append(output, fmt.Sprintf("%s=%v", tag, v.Field(i)))
	}
	return strings.Join(output, ",")
}

// From https://github.com/ceph/ceph-client/blob/for-linus/drivers/block/rbd.c#L851

//static const struct fs_parameter_spec rbd_parameters[] = {
//	fsparam_u32	("alloc_size",			Opt_alloc_size),
//	fsparam_flag	("exclusive",			Opt_exclusive),
//	fsparam_flag	("lock_on_read",		Opt_lock_on_read),
//	fsparam_u32	("lock_timeout",		Opt_lock_timeout),
//	fsparam_flag	("notrim",			Opt_notrim),
//	fsparam_string	("_pool_ns",			Opt_pool_ns),
//	fsparam_u32	("queue_depth",			Opt_queue_depth),
//	fsparam_flag	("read_only",			Opt_read_only),
//	fsparam_flag	("read_write",			Opt_read_write),
//	fsparam_flag	("ro",				Opt_read_only),
//	fsparam_flag	("rw",				Opt_read_write),
//	{}
//};

//struct rbd_options {
//	int	queue_depth;
//	int	alloc_size;
//	unsigned long	lock_timeout;
//	bool	read_only;
//	bool	lock_on_read;
//	bool	exclusive;
//	bool	trim;
//};
//
//#define RBD_QUEUE_DEPTH_DEFAULT	BLKDEV_MAX_RQ
//#define RBD_ALLOC_SIZE_DEFAULT	(64 * 1024)
//#define RBD_LOCK_TIMEOUT_DEFAULT 0  /* no timeout */
//#define RBD_READ_ONLY_DEFAULT	false
//#define RBD_LOCK_ON_READ_DEFAULT false
//#define RBD_EXCLUSIVE_DEFAULT	false
//#define RBD_TRIM_DEFAULT	true
