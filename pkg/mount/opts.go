// Majority of the below code is from https://github.com/u-root/u-root/blob/master/cmds/core/mount/opts.go
// adapted into a package. That code has the below license and copyright.
//
// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux

package mount

var opts map[string]uintptr
