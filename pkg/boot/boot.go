package boot

import (
	"log"
	"os"
)

// UnshareRoot spawns the image in newRoot as a new container.
// This is a blocking execution.
// TODO: add process tracking capabilities, non-blocking execution
// UnshareRoot is a wrapper for OS specific implementation(s)
func UnshareRoot(newRoot, init string) error {
	return unshareRoot(newRoot, init)
}

// SwitchRoot implements a switch_root to the image in newRoot and Execs init
// Since this uses Exec, the process is completely taken over by this call, so don't expect a return.
// SwitchRoot is a wrapper for OS specific implementation(s)
func SwitchRoot(newRoot, init string) error {
	return switchRoot(newRoot, init)
}

// PIDInit initializes a PID namespace and execs the provided init
// It can be called at the begining of any main() function and will only do anythin
// if the environment indicates that it should.
func PIDInit() {
	init := os.Getenv("RBD_INIT")
	if init == "" {
		return
	}
	pidInit(init)
	os.Exit(0) // need to stop execution should this return
}

// logging should probably have a global context, but this will simplify things for now
func init() {
	log.SetPrefix("boot: ")
}
