/* uinit.go: a simple init launcher for kraken layer0
 *
 * Author: J. Lowell Wofford <lowell@lanl.gov>
 *
 * This software is open source software available under the BSD-3 license.
 * Copyright (c) 2018, Triad National Security, LLC
 * See LICENSE file for details.
 */

package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

// info about background processes, inclding std{out,in,err} will be placed here
const ioDir = "/tmp/io"

// perm mask for ioDir files
const ioMode = 0700

// log level of child node
const logLevel = "7"

type command struct {
	Cmd        string
	Args       []string
	Background bool
	Exec       bool
}

const (
	kernArgFile = "/proc/cmdline"
)

func goExec(cmd *exec.Cmd) {
	if err := cmd.Run(); err != nil {
		log.Printf("command %s failed: %s\n", cmd.Path, err)
	}
}

func main() {
	log.Println("starting uinit")

	// This defines what will run
	var cmdList = []command{
		// give the system 2 seconds to come to its senses
		// shouldn't be necessary, but seems to help
		{
			Cmd:  "/bbin/sleep",
			Args: []string{"/bbin/sleep", "2"},
		},
		{
			Cmd:  "/bbin/modprobe",
			Args: []string{"modprobe", "-a", "virtio_net", "virtio-rng", "virtio_blk", "af_packet", "rbd", "squashfs", "overlay"},
		},
		{
			Cmd:  "/bbin/dhclient",
			Args: []string{"/bbin/dhclient", "-ipv6=false", "eth0"},
		},
		{
			Cmd:  "/bbin/rbd",
			Args: []string{"/bbin/rbd", "--verbose", "boot", "--mkdir", "--switch-root"},
			Exec: true,
		},
	}

	envs := os.Environ()
	os.MkdirAll(ioDir, os.ModeDir&ioMode)
	var cmds []*exec.Cmd
	for i, c := range cmdList {
		v := c.Cmd
		if _, err := os.Stat(v); !os.IsNotExist(err) {
			cmd := exec.Command(v)
			cmd.Dir = "/"
			cmd.Env = envs
			log.Printf("cmd: %d/%d %s", i+1, len(cmdList), cmd.Path)
			if c.Background {
				cmdIODir := ioDir + "/" + strconv.Itoa(i)
				if err = os.Mkdir(cmdIODir, os.ModeDir&ioMode); err != nil {
					log.Println(err)
				}
				if err = ioutil.WriteFile(cmdIODir+"/"+"cmd", []byte(v), os.ModePerm&ioMode); err != nil {
					log.Println(err)
				}
				stdin, err := os.OpenFile(cmdIODir+"/"+"stdin", os.O_RDONLY|os.O_CREATE, ioMode&os.ModePerm)
				if err != nil {
					log.Println(err)
				}
				stdout, err := os.OpenFile(cmdIODir+"/"+"stdout", os.O_WRONLY|os.O_CREATE, ioMode&os.ModePerm)
				if err != nil {
					log.Println(err)
				}
				stderr, err := os.OpenFile(cmdIODir+"/"+"stderr", os.O_WRONLY|os.O_CREATE, ioMode&os.ModePerm)
				if err != nil {
					log.Println(err)
				}
				defer stdin.Close()
				defer stdout.Close()
				defer stderr.Close()
				cmd.Stdin = stdin
				cmd.Stdout = stdout
				cmd.Stderr = stderr
				cmd.Args = c.Args
				go goExec(cmd)
				cmds = append(cmds, cmd)
			} else if c.Exec {
				if err := syscall.Exec(c.Cmd, c.Args, os.Environ()); err != nil {
					log.Printf("command %s failed: %s", cmd.Path, err)
				}
			} else {
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
				cmd.Args = c.Args
				if err := cmd.Run(); err != nil {
					log.Printf("command %s failed: %s", cmd.Path, err)
				}
			}
		}
	}
	log.Println("uinit exit")
}
