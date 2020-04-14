package cmdline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/bensallen/rbdmap/krbd"
)

type Mount struct {
	Image     *krbd.Image
	MountOpts string `json:"mntopts"`
	Part      string
	Overlay   bool
	Path      string
	Fstype    string
}

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
// or more Images. Only JSON formats below is implemented.
//
// rbd.<name>... where <name> is an arbitrary string identifer for the mount
// rbd.root.image=test-image1
// rbd.root.image.pool=rbd
// rbd.root.image.mons=192.168.0.1,192.168.0.2,192.168.0.3:6789
// rbd.root.image.user=admin
// rbd.root.image.secret=<key>
//
// Optional
// rbd.root.image.snap=snap1
// rbd.root.image.opts=rw,share
// rbd.root.part=1
// rbd.root.mntopts=defaults
// rbd.root.fstype=ext4
// rbd.root.overlay=false
// rbd.root.path=/newroot
//
// JSON
// rbd={"root": {"image":{"mons": ["192.168.0.1","192.168.0.2","192.168.0.3:6789"], "opts":{"name": "admin", "secret": "AQAvjX9eabfZAhAAj/g5nXSe/uaemYGCu1w53Q=="}, "pool":"rbd", "image":"test-image1"}, "path":"/newroot", "fstype":"ext4"}}
// rbd.root={"image":{"mons": ["192.168.0.1","192.168.0.2","192.168.0.3:6789"], "opts":{"name": "admin", "secret": "AQAvjX9eabfZAhAAj/g5nXSe/uaemYGCu1w53Q=="}, "pool":"rbd", "image":"test-image1"}, "path":"/newroot", "fstype":"ext4"}
func Parse(cmdline string) map[string]*Mount {
	log.Printf("Debug: %s", cmdline)

	mounts := map[string]*Mount{}
	for _, part := range split(cmdline) {
		switch {
		case strings.HasPrefix(part, prefix+"."):
			splitN := strings.IndexRune(part, '=')
			if splitN > 0 {
				keySplit := strings.Split(part[len(prefix)+1:splitN], ".")
				mount := &Mount{}
				// Check if a mount is already in the map
				if _, ok := mounts[keySplit[0]]; ok {
					mount = mounts[keySplit[0]]
				}

				switch len(keySplit) {
				case 1:
					// Image label and no attribute as part of key, eg. rbd.root=
					// so assume the value is JSON.
					if err := json.Unmarshal([]byte(part[splitN+1:]), mount); err != nil {
						//log.Printf("Error parsing json: %v", err)
						continue
					}
				case 2:
					// Volume label with attribute, eg. rbd.root.image=
					// [TODO]
				default:
					continue
				}
				mounts[keySplit[0]] = mount
			}
		case strings.HasPrefix(part, prefix+"="):
			// Bare rbd key, assume value is JSON
			if err := json.Unmarshal([]byte(part[len(prefix)+1:]), &mounts); err != nil {
				log.Printf("Error parsing json: %v", err)
				continue
			}
		default:
			continue
		}
	}
	return mounts
}

//Split strings on spaces except when a space is within a quoted, bracketed, or braced string.
//Supports nesting multiple brackets or braces.
func split(s string) []string {
	lastRune := map[rune]int{}
	f := func(c rune) bool {
		switch {
		case lastRune[c] > 0:
			lastRune[c]--
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastRune[c]++
			return false
		case c == '[':
			lastRune[']']++
			return false
		case c == '{':
			lastRune['}']++
			return false
		case mapGreaterThan(lastRune, 0):
			return false
		default:
			return c == ' '
		}
	}
	return strings.FieldsFunc(s, f)
}

// mapGreaterThan ranges across the provided map[rune]int looking for any values greater than
func mapGreaterThan(runes map[rune]int, g int) bool {
	for _, i := range runes {
		if i > g {
			return true
		}
	}
	return false
}

// Read open and returns all of the contents of the provided file.
func Read(path string) ([]byte, error) {
	r, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer r.Close()
	if err != nil {
		return nil, fmt.Errorf("Error opening %s: %v", path, err)
	}
	return ioutil.ReadAll(r)
}
