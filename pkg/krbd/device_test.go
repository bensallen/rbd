package krbd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_devices(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []Device
		wantErr bool
	}{
		{
			name: "Mock Dev 0",
			args: args{path: "test/sys/bus/rbd/devices"},
			want: []Device{{ID: 0, Pool: "rbd", Namespace: "ns1", Image: "image1", Snapshot: "snapshot1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := devices(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("devices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("devices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDevice_readDeviceAttrs(t *testing.T) {
	type fields struct {
		ID        int64
		Pool      string
		Namespace string
		Image     string
		Snapshot  string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Device
		wantErr bool
	}{
		{
			name:   "Mock Dev 0",
			fields: fields{ID: 0},
			want:   Device{ID: 0, Pool: "rbd", Namespace: "ns1", Image: "image1", Snapshot: "snapshot1"},
			args:   args{path: "test/sys/bus/rbd/devices/0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{
				ID:        tt.fields.ID,
				Pool:      tt.fields.Pool,
				Namespace: tt.fields.Namespace,
				Image:     tt.fields.Image,
				Snapshot:  tt.fields.Snapshot,
			}
			if err := d.readDeviceAttrs(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Device.readDeviceAttrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(*d, tt.want) {
				t.Errorf("Device.readDeviceAttrs() = %v, want %v", *d, tt.want)
			}
		})
	}
}

func TestDevice_find(t *testing.T) {
	type fields struct {
		ID        int64
		Pool      string
		Namespace string
		Image     string
		Snapshot  string
	}
	type args struct {
		devices []Device
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Device
		wantErr bool
	}{
		{
			name:   "Match on image",
			fields: fields{Image: "test1"},
			args:   args{devices: []Device{{ID: 1, Image: "test1", Pool: "rbd"}}},
			want:   Device{ID: 1, Image: "test1", Pool: "rbd"},
		},
		{
			name:   "Match on all fields",
			fields: fields{Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap1"},
			args:   args{devices: []Device{{ID: 1, Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap1"}}},
			want:   Device{ID: 1, Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap1"},
		},
		{
			name:    "Empty Device",
			fields:  fields{},
			args:    args{devices: []Device{{ID: 1, Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap1"}}},
			wantErr: true,
		},
		{
			name:   "Match on all fields",
			fields: fields{Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap1"},
			args:   args{devices: []Device{{ID: 1, Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap1"}}},
			want:   Device{ID: 1, Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap1"},
		},
		{
			name:    "No match",
			fields:  fields{Image: "test1"},
			args:    args{devices: []Device{{ID: 1, Image: "test2", Pool: "rbd"}}},
			wantErr: true,
		},
		{
			name:   "Skip namespace, snapshot, and pool mismatches",
			fields: fields{Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap1"},
			args: args{devices: []Device{
				{ID: 1, Image: "test1", Pool: "rbd", Namespace: "ns2", Snapshot: "snap1"},
				{ID: 2, Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap2"},
				{ID: 3, Image: "test1", Pool: "other", Namespace: "ns1", Snapshot: "snap1"},
				{ID: 4, Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap1"},
			}},
			want: Device{ID: 4, Image: "test1", Pool: "rbd", Namespace: "ns1", Snapshot: "snap1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{
				ID:        tt.fields.ID,
				Pool:      tt.fields.Pool,
				Namespace: tt.fields.Namespace,
				Image:     tt.fields.Image,
				Snapshot:  tt.fields.Snapshot,
			}
			err := d.find(tt.args.devices)
			if (err != nil) != tt.wantErr {
				t.Errorf("Device.find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(*d, tt.want) {
				t.Errorf("Device.find() = %v, want %v", *d, tt.want)
			}
		})
	}
}

func useFixtureSysBus(t *testing.T) {
	t.Helper()
	orig := sysBusPath
	sysBusPath = "test/sys/bus/rbd"
	t.Cleanup(func() { sysBusPath = orig })
}

func TestDevices(t *testing.T) {
	useFixtureSysBus(t)
	got, err := Devices()
	if err != nil {
		t.Fatal(err)
	}
	want := []Device{{ID: 0, Pool: "rbd", Namespace: "ns1", Image: "image1", Snapshot: "snapshot1"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Devices() = %v, want %v", got, want)
	}
}

func TestDevice_Find(t *testing.T) {
	useFixtureSysBus(t)

	d := Device{Pool: "rbd", Image: "image1"}
	if err := d.Find(); err != nil {
		t.Fatal(err)
	}
	if d.Namespace != "ns1" || d.Snapshot != "snapshot1" {
		t.Errorf("Find() did not populate attributes: %+v", d)
	}

	missing := Device{Image: "nope"}
	if err := missing.Find(); err == nil {
		t.Error("Find() matched a nonexistent image")
	}

	sysBusPath = t.TempDir()
	if err := (&Device{Image: "image1"}).Find(); err == nil {
		t.Error("Find() succeeded without a devices directory")
	}
}

func Test_devicesErrors(t *testing.T) {
	t.Run("missing dir", func(t *testing.T) {
		if _, err := devices(filepath.Join(t.TempDir(), "nope")); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("non-numeric entry", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, "x"), 0755)
		if _, err := devices(dir); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("dangling symlink", func(t *testing.T) {
		dir := t.TempDir()
		os.Symlink(filepath.Join(dir, "gone"), filepath.Join(dir, "0"))
		if _, err := devices(dir); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("missing required attr", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, "0"), 0755)
		if _, err := devices(dir); err == nil {
			t.Error("expected error for missing pool/name")
		}
	})
	t.Run("sorted by id", func(t *testing.T) {
		dir := t.TempDir()
		for _, id := range []string{"10", "2"} {
			p := filepath.Join(dir, id)
			os.Mkdir(p, 0755)
			os.WriteFile(filepath.Join(p, "pool"), []byte("rbd\n"), 0644)
			os.WriteFile(filepath.Join(p, "name"), []byte("img"+id+"\n"), 0644)
		}
		got, err := devices(dir)
		if err != nil {
			t.Fatal(err)
		}
		want := []Device{{ID: 2, Pool: "rbd", Image: "img2"}, {ID: 10, Pool: "rbd", Image: "img10"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("devices() = %+v, want %+v", got, want)
		}
	})
}

func TestDevice_DevPath(t *testing.T) {
	if got := (&Device{ID: 12}).DevPath(); got != "/dev/rbd12" {
		t.Errorf("DevPath() = %s", got)
	}
}
