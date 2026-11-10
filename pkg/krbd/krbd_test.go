package krbd

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

// fakeSysBus points sysBusPath at a temporary directory containing the
// given control files and returns the directory.
func fakeSysBus(t *testing.T, files ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}
	orig := sysBusPath
	sysBusPath = dir
	t.Cleanup(func() { sysBusPath = orig })
	return dir
}

func TestBusWriters(t *testing.T) {
	tests := []struct {
		name    string
		open    func() (io.WriteCloser, error)
		files   []string
		want    string
		wantErr bool
	}{
		{"add prefers single major", RBDBusAddWriter, []string{"add", "add_single_major"}, "add_single_major", false},
		{"add falls back", RBDBusAddWriter, []string{"add"}, "add", false},
		{"add missing", RBDBusAddWriter, nil, "", true},
		{"remove prefers single major", RBDBusRemoveWriter, []string{"remove", "remove_single_major"}, "remove_single_major", false},
		{"remove falls back", RBDBusRemoveWriter, []string{"remove"}, "remove", false},
		{"remove missing", RBDBusRemoveWriter, nil, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeSysBus(t, tt.files...)
			wc, err := tt.open()
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			defer wc.Close()
			if got := filepath.Base(wc.(*os.File).Name()); got != tt.want {
				t.Errorf("opened %s, want %s", got, tt.want)
			}
		})
	}
}

// TestMapUnmapWrites checks the exact strings written to the sysfs control
// files through the public writer helpers.
func TestMapUnmapWrites(t *testing.T) {
	dir := fakeSysBus(t, "add_single_major", "remove_single_major")

	img := Image{
		Monitors: []string{"10.0.0.1:3300"},
		Pool:     "rbd",
		Image:    "img",
		Options:  &Options{Name: "admin", Secret: "key", MsMode: "secure"},
	}
	w, err := RBDBusAddWriter()
	if err != nil {
		t.Fatal(err)
	}
	if err := img.Map(w); err != nil {
		t.Fatal(err)
	}
	w.Close()

	w, err = RBDBusRemoveWriter()
	if err != nil {
		t.Fatal(err)
	}
	if err := (&Image{DevID: 3, Options: &Options{Force: true}}).Unmap(w); err != nil {
		t.Fatal(err)
	}
	w.Close()

	for file, want := range map[string]string{
		"add_single_major":    "10.0.0.1:3300 ms_mode=secure,name=admin,secret=key rbd img -",
		"remove_single_major": "3 force",
	} {
		got, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Errorf("%s = %q, want %q", file, got, want)
		}
	}
}
