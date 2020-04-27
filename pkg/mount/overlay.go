package mount

import (
	"os"
)

// Overlay prepares and mounts a R/W overlay on the dest path. Expects lower to already be
// mounted. Create's dest, upper, and work directories. Upper and work is required to be paths
// within the same filesystem.
func Overlay(lower string, upper string, work string, dest string) error {
	for _, dir := range []string{upper, work, dest} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return Mount("overlay", dest, "overlay", []string{"lowerdir=" + lower, "upperdir=" + upper, "workdir=" + work})
}
