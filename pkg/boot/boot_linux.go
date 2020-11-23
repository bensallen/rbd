package boot

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"
)

// we store these as a global var so that we could potentially have a way to update at runtime
var specialMounts = []string{"/dev", "/proc", "/sys", "/run"}

// unshareRoot starts init in a namespaced context (container)
// Currently we only do mount + pid namespaces
func unshareRoot(newRoot, init string) (err error) {
	log.SetPrefix(log.Prefix() + "clone: ")

	log.Print("starting clone")

	// 0. For a container, we want to be able to launch bare directory containers
	//    We implement this by bind mounting newRoot on itself.
	if !isMountpoint(newRoot) {
		if err = bindMountSelf(newRoot); err != nil {
			return fmt.Errorf("clone: could not self-bind mount bare directory: %v", err)
		}
	}

	// 1. Is our image valid?
	log.Print("validating image")
	if err = validateImage(newRoot); err != nil {
		return fmt.Errorf("clone: image validation failed: %v", err)
	}

	// 2. Is our init valid?
	log.Print("validating init")
	if err = validateInit(newRoot, init); err != nil {
		return fmt.Errorf("clone: init validationfailed: %v", err)
	}

	// 3. Create new namespaces
	log.Printf("creating namespaces")
	if err = unix.Unshare(syscall.CLONE_NEWNS | syscall.CLONE_NEWPID); err != nil {
		return fmt.Errorf("clone: failed to unshare namespaces: %v", err)
	}

	// 4. Do the root moving dance
	log.Print("preparing image")
	if err = moveRoot(newRoot); err != nil {
		return fmt.Errorf("switch_root: could not prepare image: %v", err)
	}

	// 5. Exec container
	c := exec.Command(init)
	if err = c.Run(); err != nil {
		return fmt.Errorf("clone: failed to start init: %v", err)
	}
	return
}

// SwitchRoot performs a switch root into an image
func switchRoot(newRoot, init string) (err error) {
	log.SetPrefix(log.Prefix() + "switch_root: ")
	var oldRoot int

	log.Print("starting switch root")

	// 1. Is our image valid?
	log.Print("validating image")
	if err = validateImage(newRoot); err != nil {
		return fmt.Errorf("switch_root: image validation failed: %v", err)
	}

	// 2. Is our init valid?
	log.Print("validating init")
	if err = validateInit(newRoot, init); err != nil {
		return fmt.Errorf("switch_root: init validationfailed: %v", err)
	}

	// 3. Open old root for later cleanup
	if oldRoot, err = unix.Open("/", unix.O_DIRECTORY, unix.O_RDONLY); err != nil {
		return fmt.Errorf("switch_root: could not open /: %v", err)
	}
	defer unix.Close(oldRoot)

	// 4. Do the root moving dance
	log.Print("preparing image")
	if err = moveRoot(newRoot); err != nil {
		return fmt.Errorf("switch_root: could not prepare image: %v", err)
	}

	// 5. Clean up old root (if its a ramdisk). This is best-effort only.
	log.Print("cleaning up old root")
	if isRamdisk(oldRoot) {
		recursiveDelete(oldRoot)
	}

	// 6. Exec init
	log.Print("executing init")
	if err = unix.Exec(init, []string{init}, []string{}); err != nil {
		return fmt.Errorf("switch_root: exec failed: %v", err)
	}
	return
}

func isRamdisk(fd int) bool {
	stat := &unix.Statfs_t{}
	if err := unix.Fstatfs(fd, stat); err != nil {
		// we don't return errors, but we assume we are *not* a ramdisk
		return false
	}
	if stat.Type == unix.TMPFS_MAGIC || stat.Type == unix.RAMFS_MAGIC {
		return true
	}
	return false
}

func getDev(fd int) (uint64, error) {
	stat := &unix.Stat_t{}
	if err := unix.Fstat(fd, stat); err != nil {
		return 0, err
	}
	return stat.Dev, nil
}

func isMountpointAt(parentDev uint64, fd int) bool {
	dev, err := getDev(fd)
	if err != nil {
		// note this behavior is slightly arbitrary
		return false
	}
	if dev != parentDev {
		return true
	}
	return false
}

func isMountpoint(path string) bool {
	var fd, pfd int
	var parentDev uint64
	var err error
	parent := filepath.Dir(path)
	if pfd, err = unix.Open(parent, unix.O_DIRECTORY, unix.O_RDONLY); err != nil {
		// note this behavior is slightly arbitrary
		return false
	}
	defer unix.Close(pfd)
	if parentDev, err = getDev(pfd); err != nil {
		return false
	}

	if fd, err = unix.Open(path, unix.O_DIRECTORY, unix.O_RDONLY); err != nil {
		return false
	}
	defer unix.Close(fd)
	return isMountpointAt(parentDev, fd)
}

func validateImage(newRoot string) (err error) {
	var fd int
	// Does the directory exist? Or, is it a directory?
	if fd, err = unix.Open(newRoot, unix.O_DIRECTORY, unix.O_RDONLY); err != nil {
		return fmt.Errorf("new root is not a valid directory")
	}
	unix.Close(fd)

	// Is it a mount point?
	if !isMountpoint(newRoot) {
		return fmt.Errorf("new root is not a mountpoint")
	}
	return
}

func validateInit(newRoot, init string) (err error) {
	var stat os.FileInfo
	var realInit string
	var exit func() error

	if exit, err = chroot(newRoot); err != nil {
		return fmt.Errorf("could not chroot into %s: %v", newRoot, err)
	}

	if realInit, err = filepath.EvalSymlinks(init); err != nil {
		return fmt.Errorf("init file could not be found: %v", err)
	}

	if err := exit(); err != nil {
		return fmt.Errorf("could not exit chroot: %v", err)
	}

	if stat, err = os.Stat(filepath.Join(newRoot, realInit)); err != nil {
		return fmt.Errorf("init file could not be opened: %v", err)
	}
	if !stat.Mode().IsRegular() {
		return fmt.Errorf("init does not reference a regular file: %v", err)
	}
	if stat.Mode()&0111 == 0 {
		return fmt.Errorf("init file is not executable: %v", err)
	}
	return
}

func moveMount(newRoot, mount string) (err error) {
	joined := filepath.Join(newRoot, mount)
	if !isMountpoint(mount) {
		return fmt.Errorf("original mountpoint does not exist")
	}
	if isMountpoint(joined) {
		// we *do* want to unmount at least
		unix.Unmount(mount, unix.MNT_DETACH)
		return fmt.Errorf("new mountpoint already mounted, old mount detached")
	}
	if err = unix.Mount(mount, joined, "", unix.MS_MOVE, ""); err != nil {
		// we still force an unmount
		unix.Unmount(mount, unix.MNT_FORCE)
		return fmt.Errorf("mount move failed, old mount force unmounted: %v", err)
	}
	return
}

func chroot(path string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	if err := unix.Chroot(path); err != nil {
		root.Close()
		return nil, err
	}

	if err := os.Chdir("/"); err != nil {
		root.Close()
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			return err
		}
		return unix.Chroot(".")
	}, nil
}

// this is the workhorse for all schemes
// it preforms the root-moving dance
func moveRoot(newRoot string) (err error) {
	// 1. move special mounts
	for _, mount := range specialMounts {
		if err := moveMount(newRoot, mount); err != nil {
			// this isn't fatal, but we should mention it
			log.Printf("warn: couldn't move mount %s: %v", mount, err)
		}
	}
	// 2. chdir to new root
	if err = os.Chdir(newRoot); err != nil {
		return fmt.Errorf("failed to chdir to new root: %v", err)
	}
	// 3. Move newRoot -> /
	if err = unix.Mount(newRoot, "/", "", unix.MS_MOVE, ""); err != nil {
		return fmt.Errorf("failed to move new root to /: %v", err)
	}
	// 4. chroot "."
	if _, err = chroot("."); err != nil {
		return fmt.Errorf("failed to change root: %v", err)
	}

	// the dance is done
	return
}

// note: this is best-effort, and doesn't return errors
//       this is very much by design, since stopping when we hit errors could leave us in strange states
//       this function is inspired by the one found in u-root
func recursiveDelete(fd int) {
	// fd always points to a dir
	// let's use os instead of writing our own Readdirnames
	var dev uint64
	dev, err := getDev(fd)
	if err != nil {
		// odd, but just retun
		return
	}
	dir := os.NewFile(uintptr(fd), "__ignored__")
	names, err := dir.Readdirnames(-1) // does not include . or ..
	if err != nil {
		// odd, but just return
		return
	}
	for _, name := range names {
		if cfd, err := unix.Openat(fd, name, unix.O_DIRECTORY|unix.O_NOFOLLOW, unix.O_RDONLY); err != nil {
			// this is *not* a directory
			if err := unix.Unlinkat(fd, name, 0); err != nil {
				log.Printf("warn: unable to remove file %s: %v", name, err)
			}
		} else {
			// this is a directory
			if isMountpointAt(dev, cfd) {
				// we should leave mount points alone, and not descend into them
				unix.Close(cfd)
				continue
			}
			recursiveDelete(cfd)
			// recurse done, clean up this dir
			unix.Close(cfd)
			if err := unix.Unlinkat(fd, name, 0); err != nil {
				log.Printf("warn: unable to remove dir %s: %v", name, err)
			}
		}
	}

	return
}

func bindMountSelf(path string) (err error) {
	// if we're already a mount point, just return
	if isMountpoint(path) {
		return
	}
	// we blindly try this without verifying that it's a directory
	if err = unix.Mount(path, path, "", unix.MS_BIND, ""); err != nil {
		return fmt.Errorf("failed to create root bind mount: %v", err)
	}
	return
}

func unshareMount() (err error) {
	return
}
