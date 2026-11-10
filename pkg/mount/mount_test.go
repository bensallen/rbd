//go:build linux

package mount

import (
	"testing"

	"golang.org/x/sys/unix"
)

func Test_parseOptions(t *testing.T) {
	tests := []struct {
		name      string
		options   []string
		wantFlags uintptr
		wantData  string
		wantLoop  bool
	}{
		{name: "none"},
		{name: "rw is a no-op", options: []string{"rw"}},
		{name: "ro", options: []string{"ro"}, wantFlags: unix.MS_RDONLY},
		{
			name:      "flags combine",
			options:   []string{"ro", "nosuid", "nodev", "noexec", "relatime"},
			wantFlags: unix.MS_RDONLY | unix.MS_NOSUID | unix.MS_NODEV | unix.MS_NOEXEC | unix.MS_RELATIME,
		},
		{
			name:     "unknown options become data",
			options:  []string{"noload", "errors=remount-ro"},
			wantData: "noload,errors=remount-ro",
		},
		{
			name:      "mixed flags, data, and loop",
			options:   []string{"loop", "ro", "lowerdir=/a", "noatime", "upperdir=/b"},
			wantFlags: unix.MS_RDONLY | unix.MS_NOATIME,
			wantData:  "lowerdir=/a,upperdir=/b",
			wantLoop:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags, data, loop := parseOptions(tt.options)
			if flags != tt.wantFlags {
				t.Errorf("flags = %#x, want %#x", flags, tt.wantFlags)
			}
			if data != tt.wantData {
				t.Errorf("data = %q, want %q", data, tt.wantData)
			}
			if loop != tt.wantLoop {
				t.Errorf("loop = %v, want %v", loop, tt.wantLoop)
			}
		})
	}
}

func TestMount_NoFsType(t *testing.T) {
	if err := Mount("/dev/null", t.TempDir(), "", nil); err == nil {
		t.Error("expected error without a filesystem type")
	}
}
