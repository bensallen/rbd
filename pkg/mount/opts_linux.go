// Majority of the below code is from https://github.com/u-root/u-root/blob/master/cmds/core/mount/opts_linux.go
// adapted into a package. That code has the below license and copyright.

// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount

import "golang.org/x/sys/unix"

var opts = map[string]uintptr{
	"rw":           0x0, // RW is the absence of MS_RDONLY so no-op
	"active":       unix.MS_ACTIVE,
	"async":        unix.MS_ASYNC,
	"bind":         unix.MS_BIND,
	"born":         unix.MS_BORN,
	"dirsync":      unix.MS_DIRSYNC,
	"invalidate":   unix.MS_INVALIDATE,
	"i_version":    unix.MS_I_VERSION,
	"kernmount":    unix.MS_KERNMOUNT,
	"lazytime":     unix.MS_LAZYTIME,
	"mandlock":     unix.MS_MANDLOCK,
	"mgc_msk":      unix.MS_MGC_MSK,
	"mgc_val":      unix.MS_MGC_VAL,
	"move":         unix.MS_MOVE,
	"noatime":      unix.MS_NOATIME,
	"nodev":        unix.MS_NODEV,
	"nodiratime":   unix.MS_NODIRATIME,
	"noexec":       unix.MS_NOEXEC,
	"noremotelock": unix.MS_NOREMOTELOCK,
	"nosec":        unix.MS_NOSEC,
	"nosuid":       unix.MS_NOSUID,
	//"nouser":       unix.MS_NOUSER, // Looks like this results in EINVAL and is removed in more recent kernels
	"posixacl":    unix.MS_POSIXACL,
	"private":     unix.MS_PRIVATE,
	"rdonly":      unix.MS_RDONLY,
	"ro":          unix.MS_RDONLY, // Alias to rdonly
	"rec":         unix.MS_REC,
	"relatime":    unix.MS_RELATIME,
	"remount":     unix.MS_REMOUNT,
	"rmt_mask":    unix.MS_RMT_MASK,
	"shared":      unix.MS_SHARED,
	"silent":      unix.MS_SILENT,
	"slave":       unix.MS_SLAVE,
	"strictatime": unix.MS_STRICTATIME,
	"submount":    unix.MS_SUBMOUNT,
	"sync":        unix.MS_SYNC,
	"synchronous": unix.MS_SYNCHRONOUS,
	"unbindable":  unix.MS_UNBINDABLE,
	"verbose":     unix.MS_VERBOSE,
}
