package krbd

import (
	"bytes"
	"testing"
)

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
