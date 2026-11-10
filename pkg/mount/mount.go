// Below code is adapted from https://github.com/u-root/u-root/blob/master/cmds/core/mount/mount.go
// That code has the below license and copyright.
//
// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount

import (
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/loop"
)

func loopSetup(filename string) (loopDevice string, err error) {
	loopDevice, err = loop.FindDevice()
	if err != nil {
		return "", err
	}
	if err := loop.SetFile(loopDevice, filename); err != nil {
		return "", err
	}
	return loopDevice, nil
}

// parseOptions splits mount options into mount(2) flags and a comma separated
// filesystem data string. loop reports whether the "loop" option was given.
func parseOptions(options []string) (flags uintptr, data string, loop bool) {
	var d []string
	for _, option := range options {
		if option == "loop" {
			loop = true
			continue
		}
		if f, ok := opts[option]; ok {
			flags |= f
		} else {
			d = append(d, option)
		}
	}
	return flags, strings.Join(d, ","), loop
}

// Mount extends u-root/pkg/mount to setup loop devices and parse input options
// into data and flags
func Mount(dev string, path string, fsType string, options []string) error {
	if fsType == "" {
		return fmt.Errorf("no file system type provided")
	}

	flags, data, loop := parseOptions(options)
	if loop {
		var err error
		if dev, err = loopSetup(dev); err != nil {
			return fmt.Errorf("error setting loop device: %s", err)
		}
	}

	_, err := mount.Mount(dev, path, fsType, data, flags)
	return err
}
