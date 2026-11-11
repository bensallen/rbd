//go:build smoke && linux

package smoke

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bensallen/rbd/pkg/krbd"
	"github.com/bensallen/rbd/pkg/mount"
)

const (
	pool  = "rbd"
	image = "smoke"
	user  = "smoke"
)

var payload = []byte("hello from rbd smoke test\n")

type env struct {
	mon    string
	secret string
	rbdBin string
}

func setup(t *testing.T) env {
	t.Helper()
	e := env{
		mon:    os.Getenv("RBD_SMOKE_MON"),
		secret: os.Getenv("RBD_SMOKE_SECRET"),
		rbdBin: os.Getenv("RBD_SMOKE_BIN"),
	}
	if e.mon == "" || e.secret == "" {
		t.Skip("RBD_SMOKE_MON and RBD_SMOKE_SECRET must be set")
	}
	// Past this point the environment is meant to be complete, so missing
	// prerequisites fail rather than skip.
	if os.Geteuid() != 0 {
		t.Fatal("requires root")
	}
	if _, err := os.Stat("/sys/bus/rbd"); err != nil {
		t.Fatal("rbd kernel module not loaded")
	}
	return e
}

func mapImage(t *testing.T, mons []string, opts krbd.Options) krbd.Device {
	t.Helper()
	img := krbd.Image{Monitors: mons, Pool: pool, Image: image, Options: &opts}
	w, err := krbd.RBDBusAddWriter()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	if err := img.Map(w); err != nil {
		t.Fatalf("map %s: %v", img, err)
	}
	dev := krbd.Device{Pool: pool, Image: image}
	if err := dev.Find(); err != nil {
		t.Fatalf("find: %v", err)
	}
	waitFor(t, dev.DevPath())
	return dev
}

func unmapDevice(t *testing.T, id int64) {
	t.Helper()
	w, err := krbd.RBDBusRemoveWriter()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	if err := (&krbd.Image{DevID: int(id)}).Unmap(w); err != nil {
		t.Fatalf("unmap: %v", err)
	}
	assertNoDevices(t)
}

func assertNoDevices(t *testing.T) {
	t.Helper()
	devs, err := krbd.Devices()
	if err != nil {
		t.Fatal(err)
	}
	if len(devs) != 0 {
		t.Fatalf("devices still mapped: %+v", devs)
	}
}

func waitFor(t *testing.T, path string) {
	t.Helper()
	for i := 0; i < 50; i++ {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("%s did not appear", path)
}

func readPayload(t *testing.T, path string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("content mismatch: %q", got)
	}
}

// TestLibrary maps the image with each messenger mode, formatting and writing
// on the first pass and verifying the data on subsequent passes.
func TestLibrary(t *testing.T) {
	e := setup(t)
	mnt := t.TempDir()

	cases := []struct {
		name string
		mons []string
		opts krbd.Options
	}{
		{"msgr2-prefer-crc", []string{e.mon + ":3300"}, krbd.Options{MsMode: "prefer-crc"}},
		{"msgr2-secure", []string{e.mon + ":3300"}, krbd.Options{MsMode: "secure"}},
		{"msgr1-legacy", []string{e.mon + ":6789"}, krbd.Options{}},
		{"msgr2-crc-readonly", []string{e.mon + ":3300"}, krbd.Options{MsMode: "crc", ReadOnly: true}},
	}

	for i, c := range cases {
		ok := t.Run(c.name, func(t *testing.T) {
			c.opts.Name = user
			c.opts.Secret = e.secret
			dev := mapImage(t, c.mons, c.opts)
			defer unmapDevice(t, dev.ID)

			if i == 0 {
				if out, err := exec.Command("mkfs.ext4", "-q", "-F", dev.DevPath()).CombinedOutput(); err != nil {
					t.Fatalf("mkfs.ext4: %v: %s", err, out)
				}
			}

			var mopts []string
			if c.opts.ReadOnly {
				mopts = []string{"ro", "noload"}
			}
			if err := mount.Mount(dev.DevPath(), mnt, "ext4", mopts); err != nil {
				t.Fatalf("mount: %v", err)
			}
			defer func() {
				if err := mount.Unmount(mnt, false, false); err != nil {
					t.Fatalf("unmount: %v", err)
				}
			}()

			f := filepath.Join(mnt, "hello.txt")
			if i == 0 {
				if err := os.WriteFile(f, payload, 0644); err != nil {
					t.Fatal(err)
				}
			}
			readPayload(t, f)

			if c.opts.ReadOnly {
				if err := os.WriteFile(filepath.Join(mnt, "nope"), payload, 0644); err == nil {
					t.Fatal("write succeeded on read-only mapping")
				}
			}
		})
		if !ok {
			t.FailNow()
		}
	}
}

func rbdCLI(t *testing.T, e env, args ...string) string {
	t.Helper()
	if e.rbdBin == "" {
		t.Fatal("RBD_SMOKE_BIN not set")
	}
	out, err := exec.Command(e.rbdBin, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("rbd %s: %v: %s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

func TestCLIMapListUnmap(t *testing.T) {
	e := setup(t)
	rbdCLI(t, e, "map", "-m", e.mon+":3300", "-p", pool, "-i", image, "--id", user, "--secret", e.secret, "--ms-mode", "crc")

	dev := krbd.Device{Pool: pool, Image: image}
	if err := dev.Find(); err != nil {
		t.Fatalf("find: %v", err)
	}
	if out := rbdCLI(t, e, "device", "list"); !strings.Contains(out, dev.DevPath()) {
		t.Errorf("device list missing %s:\n%s", dev.DevPath(), out)
	}
	rbdCLI(t, e, "unmap", "-d", fmt.Sprint(dev.ID))
	assertNoDevices(t)
}

// TestCLIBoot maps and mounts the image from a cmdline file with an overlay
// root, as rbd boot would in an initramfs.
func TestCLIBoot(t *testing.T) {
	e := setup(t)
	const (
		root  = "/newroot"
		lower = "/run/overlayfs/lower"
		upper = "/run/overlayfs/upper"
	)

	spec, err := json.Marshal(map[string]any{
		"image": map[string]any{
			"mons":  []string{e.mon + ":3300"},
			"pool":  pool,
			"image": image,
			"opts":  map[string]any{"name": user, "secret": e.secret, "msmode": "prefer-crc"},
		},
		"path":    "/",
		"fstype":  "ext4",
		"overlay": true,
	})
	if err != nil {
		t.Fatal(err)
	}
	cmdline := filepath.Join(t.TempDir(), "cmdline")
	if err := os.WriteFile(cmdline, []byte("console=ttyS0 rbd.root="+string(spec)), 0644); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		mount.Unmount(root, false, false)
		mount.Unmount(lower, false, false)
		dev := krbd.Device{Pool: pool, Image: image}
		if dev.Find() == nil {
			unmapDevice(t, dev.ID)
		}
	})

	rbdCLI(t, e, "boot", "--mkdir", "--cmdline", cmdline)

	readPayload(t, filepath.Join(root, "hello.txt"))
	if err := os.WriteFile(filepath.Join(root, "upper.txt"), payload, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(upper, "upper.txt")); err != nil {
		t.Errorf("write did not land in overlay upper dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(lower, "upper.txt")); err == nil {
		t.Error("write leaked into lower rbd filesystem")
	}
}
