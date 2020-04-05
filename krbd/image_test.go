package krbd

import (
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
			if got := i.String(); got != tt.want {
				t.Errorf("Image.String() = %v want %v", got, tt.want)
			}
		})
	}
}
