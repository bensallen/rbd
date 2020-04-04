package krbd

import (
	"fmt"
	"testing"
)

func TestImage_String(t *testing.T) {
	type fields struct {
		DevID     int
		Monitors  []string
		Pool      string
		Namespace string
		Image     string
		Options   *Options
		Snapshot  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Three Monitors, Typical Options",
			fields: fields{
				Monitors: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
				Pool:     "rbd",
				Image:    "test-image",
				Options:  &Options{Name: "admin", Secret: "AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w=="},
			},
			want: "10.0.0.1,10.0.0.2,10.0.0.3 name=admin,secret=AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w== rbd test-image -",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Image{
				DevID:     tt.fields.DevID,
				Monitors:  tt.fields.Monitors,
				Pool:      tt.fields.Pool,
				Namespace: tt.fields.Namespace,
				Image:     tt.fields.Image,
				Options:   tt.fields.Options,
				Snapshot:  tt.fields.Snapshot,
			}
			if got := fmt.Sprintf("%s", i); got != tt.want {
				t.Errorf("Image.String() via  fmt.Sprintf(\"%%s\", i) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptions_String(t *testing.T) {
	type fields struct {
		Exclusive   bool
		LockOnRead  bool
		NoTrim      bool
		ReadOnly    bool
		AllocSize   int
		QueueDepth  int
		LockTimeout uint64
		Name        string
		Namespace   string
		Secret      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Empty",
			fields: fields{},
			want:   "",
		},
		{
			name:   "Name and Secret",
			fields: fields{Name: "admin", Secret: "AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w=="},
			want:   "name=admin,secret=AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w==",
		},
		{
			name:   "All fields",
			fields: fields{Exclusive: true, LockOnRead: true, NoTrim: true, ReadOnly: true, AllocSize: 65536, QueueDepth: 128, LockTimeout: 500, Name: "admin", Namespace: "ns1", Secret: "AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w=="},
			want:   "exclusive=true,lock_on_read=true,notrim=true,read_only=true,alloc_size=65536,queue_depth=128,lock_timeout=500,name=admin,_pool_ns=ns1,secret=AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{
				Exclusive:   tt.fields.Exclusive,
				LockOnRead:  tt.fields.LockOnRead,
				NoTrim:      tt.fields.NoTrim,
				ReadOnly:    tt.fields.ReadOnly,
				AllocSize:   tt.fields.AllocSize,
				QueueDepth:  tt.fields.QueueDepth,
				LockTimeout: tt.fields.LockTimeout,
				Name:        tt.fields.Name,
				Namespace:   tt.fields.Namespace,
				Secret:      tt.fields.Secret,
			}
			if got := fmt.Sprintf("%s", o); got != tt.want {
				t.Errorf("Options.String() via fmt.Sprintf(\"%%s\", o) = %v, want %v", got, tt.want)
			}
		})
	}
}
