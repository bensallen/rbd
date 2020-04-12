package cmdline

import (
	"reflect"
	"testing"
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
			name: "Strings with JSON Slice",
			args: args{s: "test1 test2=[ '192.168.64.1' ]"},
			want: []string{"test1", "test2=[ '192.168.64.1' ]"},
		},
		{
			name: "Strings with JSON Map",
			args: args{s: "test1 test2={'test3': '192.168.64.1'}"},
			want: []string{"test1", "test2={'test3': '192.168.64.1'}"},
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
		want map[string]*Image
	}{
		{
			name: "rbd.root=",
			args: args{cmdline: `rbd.root={"pool":"rbd", "image":"test-image1", "path":"/newroot", "fstype":"ext4"}`},
			want: map[string]*Image{"root": {Pool: "rbd", Image: "test-image1", Path: "/newroot", Fstype: "ext4"}},
		},
		{
			name: "rbd.root= specified twice with different attributes",
			args: args{cmdline: `rbd.root={"pool":"rbd", "image":"test-image1"} rbd.root={"path":"/newroot", "fstype":"ext4"}`},
			want: map[string]*Image{"root": {Pool: "rbd", Image: "test-image1", Path: "/newroot", Fstype: "ext4"}},
		},
		{
			name: "rbd=",
			args: args{cmdline: `rbd={"root": {"pool":"rbd", "image":"test-image1", "path":"/newroot", "fstype":"ext4"}}`},
			want: map[string]*Image{"root": {Pool: "rbd", Image: "test-image1", Path: "/newroot", Fstype: "ext4"}},
		},
		{
			name: "Garbage JSON",
			args: args{cmdline: `rbd={root: asdf}}`},
			want: map[string]*Image{},
		},
		{
			name: "Garbage JSON 2",
			args: args{cmdline: `rbd.root={root: asdf}}`},
			want: map[string]*Image{},
		},
		{
			name: "Malformed key",
			args: args{cmdline: "rbd.root.pool.test=pool1"},
			want: map[string]*Image{},
		},
		{
			name: "Unrelated cmdline args no rbd",
			args: args{cmdline: "root=/dev/mapper/vg0-lv_root ro modules=sd-mod,usb-storage,ext4 nomodeset earlyprintk=ttyS0 console=ttyS0 rootfstype=ext4"},
			want: map[string]*Image{},
		},
		{
			name: "Unrelated cmdline args with rbd",
			args: args{cmdline: `nomodeset earlyprintk=ttyS0 console=ttyS0 rbd={"root": {"pool":"rbd", "image":"test-image1", "path":"/newroot", "fstype":"ext4"}}`},
			want: map[string]*Image{"root": {Pool: "rbd", Image: "test-image1", Path: "/newroot", Fstype: "ext4"}},
		},

		//
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
