package cmdline

import (
	"encoding/json"
	"log"
	"strings"
	"unicode"
)

// Image is a RBD image, which was generated via parsed data from the kernel cmdline.
type Image struct {
	DevID     int
	Monitors  []string `json:"mon"`
	Pool      string
	Image     string
	Snapshot  string `json:"snap"`
	Options   string `json:"opts"`
	MountOpts string `json:"mntopts"`
	Part      string
	Overlay   bool
	Path      string
	Fstype    string
}

// Leading prefix for cmdline arguments
const prefix = "rbd"

// Parse attempts to find rbd options from input kernel cmdline and return one
// or more Images
//
// rbd.<name>... where <name> is an arbitrary string identifer
// rbd.root.pool=rbd
// rbd.root.image=test-image1
// rbd.root.mon=192.168.0.1,192.168.0.2,192.168.0.3:6789
// rbd.root.user=admin
// rbd.root.secret=<key>
//
// Optional
// rbd.root.snapshot=snap1
// rbd.root.part=1
// rbd.root.opts=rw,share
// rbd.root.mntopts=defaults
// rbd.root.fstype=ext4
// rbd.root.overlay=false
// rbd.root.path=/newroot
//
// JSON
// rbd={"root": {"pool":"rbd", "image":"test-image1", "path":"/newroot", "fstype":"ext4"}}
// rbd.root={"pool":"rbd", "image":"test-image1", "path":"/newroot", "fstype":"ext4"}
//
func Parse(cmdline string) map[string]*Image {
	images := map[string]*Image{}
	for _, part := range split(cmdline) {
		switch {
		case strings.HasPrefix(part, prefix+"."):
			splitN := strings.IndexRune(part, '=')
			if splitN > 0 {
				keySplit := strings.Split(part[len(prefix)+1:splitN], ".")
				image := &Image{}
				// Check if a vol is already in the map
				if _, ok := images[keySplit[0]]; ok {
					image = images[keySplit[0]]
				}

				switch len(keySplit) {
				case 1:
					// Image label and no attribute as part of key, eg. rbd.root=
					// so assume the value is JSON.
					if err := json.Unmarshal([]byte(part[splitN+1:]), image); err != nil {
						log.Printf("Error parsing json: %v", err)
						continue
					}
				case 2:
					// Volume label with attribute, eg. rbd.root.image=
					// TODO Unmarshaler that matches keySplit[1] to JSON tag of attributes in Image
				default:
					continue
				}
				images[keySplit[0]] = image
			}
		case strings.HasPrefix(part, prefix+"="):
			// Bare rbd key, assume value is JSON
			if err := json.Unmarshal([]byte(part[len(prefix)+1:]), &images); err != nil {
				log.Printf("Error parsing json: %v", err)
				continue
			}
		default:
			continue
		}
	}
	return images
}

//split splits on spaces except when a space is within a quoted or bracketed string.
func split(s string) []string {
	lastRune := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastRune:
			lastRune = rune(0)
			return false
		case lastRune != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastRune = c
			return false
		case c == '[':
			lastRune = ']'
			return false
		case c == '{':
			lastRune = '}'
			return false
		default:
			return c == ' '
		}
	}
	return strings.FieldsFunc(s, f)
}
