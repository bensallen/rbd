package krbd

import (
	"testing"
)

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
			if got := o.String(); got != tt.want {
				t.Errorf("Options.String() = %v want %v", got, tt.want)
			}
		})
	}
}
