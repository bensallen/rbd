package mount

import (
	"os"
)

// Overlay prepares and mounts a R/W overlay over the provided path. Mounts /run as a tmpfs and uses
// /run/overlayfs/rw as the upper and /run/overlayfs/work as the workdir.
func Overlay(path string) error {
	if err := os.MkdirAll("/run", 0755); err != nil {
		return err
	}

	if err := Mount("tmpfs", "/run", "tmpfs", []string{"rw", "nosuid", "nodev", "mode=755"}); err != nil {
		return err
	}

	if err := os.MkdirAll("/run/overlayfs/rw", 0755); err != nil {
		return err
	}
	if err := os.MkdirAll("/run/overlayfs/work", 0755); err != nil {
		return err
	}

	return Mount("overlay", "/newroot", "overlay", []string{"lowerdir=" + path, "upperdir=/run/overlayfs/rw", "workdir=/run/overlayfs/work"})
}
