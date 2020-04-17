# Go RBD CLI

Golang based Ceph RBD image map, unmap, boot, and device list utility and library. No external dependencies outside of the Kernel's RBD sysfs interface.

```
$ rbd -h
rbd - Ceph RBD CLI

Usage:
  rbd [map|unmap|device|boot]

Subcommands:
  map      Map RBD Image
  unmap    Unmap RBD Image
  boot     Boot via RBD Image
  device   Manage RBD Devices

Flags:
  -h, --help      Diplay help.
  -n, --noop      No-op (don't actually perform action).
  -v, --verbose   Enable additional output.
  -V, --version   Displays the program version string.
```

## map

```
$ rbd map -h
map - Map RBD Image

Usage:
  map

Flags:
      --id string          Specifies the username (without the 'client.' prefix)
  -i, --image string       Image to map
  -m, --monitor strings    Connect to one or more monitor addresses (192.168.0.1[:6789]). Multiple address are specified comma separated.
      --namespace string   Use a pre-defined image namespace within a pool
  -p, --pool string        Interact with the given pool.
      --read-only          Map the image read-only
      --secret string      Specifies the user authentication secret
      --snap string        Specifies a snapshot name
```

## unmap

```
$ rbd unmap -h
unmap - Unmap RBD Image

Usage:
  unmap

Flags:
  -d, --devid int   RBD Device Index (default 0)
  -f, --force       Optional force argument will wait for running requests and then unmap the image
```

## Boot

- Unfinished, only maps devices so far.
- Parse /proc/cmdline for RBD settings and attempt to pivot_root.
- See https://github.com/bensallen/rbd/blob/master/pkg/cmdline/cmdline.go#L50 for cmdline format

```
$ rbd boot -h
boot - Boot via RBD image

Usage:
  boot

Flags:
  -c, --cmdline string   Path to kernel cmdline (default: /proc/cmdline) (default "/proc/cmdline")
```

## Device 

```
$ rbd device -h
device - Manage RBD Devices

Usage:
  device [list|map|unmap]

Subcommands:
  list     List connected devices
  map      Map RBD Image
  unmap    Unmap RBD Image
```