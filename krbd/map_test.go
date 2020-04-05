package krbd

import (
	"bytes"
	"fmt"
	"testing"
)

func TestImage_String(t *testing.T) {
	type fields struct {
		DevID    int
		Monitors []string
		Pool     string
		Image    string
		Snapshot string
		Options  *Options
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
		{
			name:   "Empty",
			fields: fields{},
			want:   " <nil>   -",
		},
		{
			name:   "Empty with Empty Option",
			fields: fields{Options: &Options{}},
			want:   "    -",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Image{
				DevID:    tt.fields.DevID,
				Monitors: tt.fields.Monitors,
				Pool:     tt.fields.Pool,
				Image:    tt.fields.Image,
				Snapshot: tt.fields.Snapshot,
				Options:  tt.fields.Options,
			}
			if got := fmt.Sprintf("%s", i); got != tt.want {
				t.Errorf("Image.String() via  fmt.Sprintf(\"%%s\", i) = %v want %v", got, tt.want)
			}
		})
	}
}

func TestOptions_String(t *testing.T) {
	type fields struct {
		// Client Options
		Fsid                     string
		IP                       string
		Share                    bool
		Noshare                  bool
		CRC                      bool
		NoCRC                    bool
		CephxRequireSignatures   bool
		NoCephxRequireSignatures bool
		TCPNoDelay               bool
		NoTCPNoDelay             bool
		CephxSignMessages        bool
		NoCephxSignMessages      bool
		MountTimeout             int
		OSDKeepAlive             int
		OSDIdleTTL               int

		// RBD Block Options
		ReadWrite   bool
		ReadOnly    bool
		QueueDepth  int
		LockOnRead  bool
		Exclusive   bool
		LockTimeout uint64
		NoTrim      bool
		AbortOnFull bool
		AllocSize   int
		Name        string
		Secret      string
		Namespace   string
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
			name: "All RBD fields",
			fields: fields{
				ReadWrite:   true,
				ReadOnly:    true,
				QueueDepth:  128,
				LockOnRead:  true,
				Exclusive:   true,
				LockTimeout: 500,
				NoTrim:      true,
				AbortOnFull: true,
				AllocSize:   65536,
				Name:        "admin",
				Secret:      "AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w==",
				Namespace:   "ns1",
			},
			want: "rw,ro,queue_depth=128,lock_on_read,exclusive,lock_timeout=500,notrim,abort_on_full,alloc_size=65536,name=admin,secret=AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w==,_pool_ns=ns1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{
				Fsid:                     tt.fields.Fsid,
				IP:                       tt.fields.IP,
				Share:                    tt.fields.Share,
				Noshare:                  tt.fields.Noshare,
				CRC:                      tt.fields.CRC,
				NoCRC:                    tt.fields.NoCRC,
				CephxRequireSignatures:   tt.fields.CephxRequireSignatures,
				NoCephxRequireSignatures: tt.fields.NoCephxRequireSignatures,
				TCPNoDelay:               tt.fields.TCPNoDelay,
				NoTCPNoDelay:             tt.fields.NoTCPNoDelay,
				CephxSignMessages:        tt.fields.CephxSignMessages,
				NoCephxSignMessages:      tt.fields.NoCephxSignMessages,
				MountTimeout:             tt.fields.MountTimeout,
				OSDKeepAlive:             tt.fields.OSDKeepAlive,
				OSDIdleTTL:               tt.fields.OSDIdleTTL,

				// RBD Block Options
				ReadWrite:   tt.fields.ReadWrite,
				ReadOnly:    tt.fields.ReadOnly,
				QueueDepth:  tt.fields.QueueDepth,
				LockOnRead:  tt.fields.LockOnRead,
				Exclusive:   tt.fields.Exclusive,
				LockTimeout: tt.fields.LockTimeout,
				NoTrim:      tt.fields.NoTrim,
				AbortOnFull: tt.fields.AbortOnFull,
				AllocSize:   tt.fields.AllocSize,
				Name:        tt.fields.Name,
				Secret:      tt.fields.Secret,
				Namespace:   tt.fields.Namespace,
			}
			if got := fmt.Sprintf("%s", o); got != tt.want {
				t.Errorf("Options.String() via fmt.Sprintf(\"%%s\", o) = %v want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Map(t *testing.T) {
	type fields struct {
		Monitors []string
		Pool     string
		Image    string
		Options  *Options
		Snapshot string
	}
	tests := []struct {
		name    string
		fields  fields
		wantW   string
		wantErr bool
	}{
		{
			name: "Typical case",
			fields: fields{
				Monitors: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
				Pool:     "rbd",
				Image:    "test-image",
				Options:  &Options{Name: "admin", Secret: "AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w=="},
			},
			wantW:   "10.0.0.1,10.0.0.2,10.0.0.3 name=admin,secret=AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w== rbd test-image -",
			wantErr: false,
		},
		{
			name: "Missing Monitor",
			fields: fields{
				Monitors: []string{},
				Pool:     "rbd",
				Image:    "test-image",
				Options:  &Options{Name: "admin", Secret: "AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w=="},
			},
			wantErr: true,
		},
		{
			name: "Missing Pool",
			fields: fields{
				Monitors: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
				Pool:     "",
				Image:    "test-image",
				Options:  &Options{Name: "admin", Secret: "AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w=="},
			},
			wantErr: true,
		},
		{
			name: "Missing Image",
			fields: fields{
				Monitors: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
				Pool:     "rbd",
				Image:    "",
				Options:  &Options{Name: "admin", Secret: "AQCvCbtToC6MDhAATtuT70Sl+DymPCfDSsyV4w=="},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Image{
				Monitors: tt.fields.Monitors,
				Pool:     tt.fields.Pool,
				Image:    tt.fields.Image,
				Options:  tt.fields.Options,
				Snapshot: tt.fields.Snapshot,
			}
			w := &bytes.Buffer{}
			if err := i.Map(w); (err != nil) != tt.wantErr {
				t.Errorf("Image.Map() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Image.Map() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestImage_Unmap(t *testing.T) {
	type fields struct {
		DevID   int
		Options *Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantW   string
		wantErr bool
	}{
		{
			name:    "No DevID (defaults to 0)",
			fields:  fields{},
			wantW:   "0",
			wantErr: false,
		},
		{
			name:    "DevID 1",
			fields:  fields{DevID: 1},
			wantW:   "1",
			wantErr: false,
		},
		{
			name:    "DevID 1 Force",
			fields:  fields{DevID: 1, Options: &Options{Force: true}},
			wantW:   "1 force",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Image{
				DevID:   tt.fields.DevID,
				Options: tt.fields.Options,
			}
			w := &bytes.Buffer{}
			if err := i.Unmap(w); (err != nil) != tt.wantErr {
				t.Errorf("Image.Unmap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Image.Unmap() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
