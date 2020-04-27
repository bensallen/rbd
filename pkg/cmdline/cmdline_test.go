package cmdline

import (
	"reflect"
	"testing"

	"github.com/bensallen/rbd/pkg/krbd"
)

func Test_split(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty",
			args: args{},
			want: []string{},
		},
		{
			name: "Strings",
			args: args{s: "root=/dev/mapper/vg0-lv_root ro modules=sd-mod,usb-storage,ext4 nomodeset earlyprintk=ttyS0 console=ttyS0 rootfstype=ext4"},
			want: []string{"root=/dev/mapper/vg0-lv_root", "ro", "modules=sd-mod,usb-storage,ext4", "nomodeset", "earlyprintk=ttyS0", "console=ttyS0", "rootfstype=ext4"},
		},
		{
			name: "Strings with JSON",
			args: args{s: "test1=test2 test3=[192.168.64.1]"},
			want: []string{"test1=test2", "test3=[192.168.64.1]"},
		},
		{
			name: "Strings with nested quotes",
			args: args{s: `"test1=test2 'test3 test4'"`},
			want: []string{`"test1=test2 'test3 test4'"`},
		},
		{
			name: "Strings with JSON Slice",
			args: args{s: "test1 test2=[ '192.168.64.1' ]"},
			want: []string{"test1", "test2=[ '192.168.64.1' ]"},
		},
		{
			name: "Strings with Nested JSON Slice",
			args: args{s: "test1 test2=[ [ '192.168.64.1' ] ]"},
			want: []string{"test1", "test2=[ [ '192.168.64.1' ] ]"},
		},
		{
			name: "Strings with JSON Map",
			args: args{s: "test1 test2={'test3': '192.168.64.1'}"},
			want: []string{"test1", "test2={'test3': '192.168.64.1'}"},
		},
		{
			name: "Strings with Nested JSON Map",
			args: args{s: "test1 test2={'test3':{ 'ip': '192.168.64.1'}}"},
			want: []string{"test1", "test2={'test3':{ 'ip': '192.168.64.1'}}"},
		},
		{
			name: "Strings with Single Quoted String",
			args: args{s: "test1 'test2 test3'"},
			want: []string{"test1", "'test2 test3'"},
		},
		{
			name: "Strings with Double Quoted String",
			args: args{s: "test1 \"test2 test3\""},
			want: []string{"test1", "\"test2 test3\""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := split(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		cmdline string
	}
	tests := []struct {
		name string
		args args
		want map[string]*Mount
	}{
		{
			name: "rbd.root=",
			args: args{cmdline: `rbd={"root": {"image":{"mons": ["192.168.0.1","192.168.0.2","192.168.0.3:6789"], "opts":{"name": "admin", "secret": "AQAvjX9eabfZAhAAj/g5nXSe/uaemYGCu1w53Q==", "readonly": true}, "pool":"rbd", "image":"test-image1"}, "path":"/newroot", "fstype":"ext4", "overlay": true}}`},
			want: map[string]*Mount{"root": {Image: &krbd.Image{Monitors: []string{"192.168.0.1", "192.168.0.2", "192.168.0.3:6789"}, Options: &krbd.Options{Name: "admin", Secret: "AQAvjX9eabfZAhAAj/g5nXSe/uaemYGCu1w53Q==", ReadOnly: true}, Pool: "rbd", Image: "test-image1"}, Path: "/newroot", FsType: "ext4", Overlay: true}},
		},
		{
			name: "rbd.root= specified twice with different attributes",
			args: args{cmdline: `rbd.root={"image":{"pool":"rbd", "image":"test-image1"}} rbd.root={"path":"/newroot", "fstype":"ext4"}`},
			want: map[string]*Mount{"root": {Image: &krbd.Image{Pool: "rbd", Image: "test-image1"}, Path: "/newroot", FsType: "ext4"}},
		},
		{
			name: "rbd=",
			args: args{cmdline: `rbd={"root":{"image":{"pool":"rbd", "image":"test-image1"}, "path":"/newroot", "fstype":"ext4"}}`},
			want: map[string]*Mount{"root": {Image: &krbd.Image{Pool: "rbd", Image: "test-image1"}, Path: "/newroot", FsType: "ext4"}},
		},
		{
			name: "Garbage JSON",
			args: args{cmdline: `rbd={"root": "asdf"}}`},
			want: map[string]*Mount{},
		},
		{
			name: "Garbage JSON 2",
			args: args{cmdline: `rbd.root={"root": "asdf"}}`},
			want: map[string]*Mount{},
		},
		{
			name: "Malformed key",
			args: args{cmdline: "rbd.root.pool.test=pool1"},
			want: map[string]*Mount{},
		},
		{
			name: "Unrelated cmdline args no rbd",
			args: args{cmdline: "root=/dev/mapper/vg0-lv_root ro modules=sd-mod,usb-storage,ext4 nomodeset earlyprintk=ttyS0 console=ttyS0 rootfstype=ext4"},
			want: map[string]*Mount{},
		},
		{
			name: "Unrelated cmdline args with rbd",
			args: args{cmdline: `nomodeset earlyprintk=ttyS0 console=ttyS0 rbd={"root":{"image":{"pool":"rbd", "image":"test-image1"}, "path":"/newroot", "fstype":"ext4"}}`},
			want: map[string]*Mount{"root": {Image: &krbd.Image{Pool: "rbd", Image: "test-image1"}, Path: "/newroot", FsType: "ext4"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.args.cmdline)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
