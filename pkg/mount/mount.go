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

	"github.com/u-root/u-root/pkg/loop"
	"github.com/u-root/u-root/pkg/mount"
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

// Mount extends u-root/pkg/mount to setup loop devices and parse input options
// into data and flags
func Mount(dev string, path string, fsType string, options []string) error {
	var err error
	var flags uintptr
	var data []string

	if fsType == "" {
		return fmt.Errorf("no file system type provided")
	}

	for _, option := range options {
		switch option {
		case "loop":
			dev, err = loopSetup(dev)
			if err != nil {
				return fmt.Errorf("error setting loop device: %s", err)
			}
		default:
			if f, ok := opts[option]; ok {
				flags |= f
			} else {
				data = append(data, option)
			}
		}
	}

	return mount.Mount(dev, path, fsType, strings.Join(data, ","), flags)
}
