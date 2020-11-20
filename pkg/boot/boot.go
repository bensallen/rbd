package boot

import "log"

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

// logging should probably have a global context, but this will simplify things for now
func init() {
	log.SetPrefix("boot: ")
}
