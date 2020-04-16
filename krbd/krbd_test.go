package krbd

import (
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("devices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDevice_readDeviceAttrs(t *testing.T) {
	type fields struct {
		ID        int
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
			}
			if !reflect.DeepEqual(*d, tt.want) {
				t.Errorf("Device.readDeviceAttrs() = %v, want %v", *d, tt.want)
			}
		})
	}
}
