package mount

import "github.com/u-root/u-root/pkg/mount"

// SwitchRoot makes newRootDir the new root directory of the system.
// Simply wraps u-root/pkg/mount.SwitchRoot
func SwitchRoot(newRootDir string, init string) error {
	return mount.SwitchRoot(newRootDir, init)
}
