package krbd

import (
	"bytes"
	"testing"
)

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
