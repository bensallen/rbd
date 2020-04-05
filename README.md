# RBD Map and Unmap

Golang based Ceph RBD image map and unmap utility and library.

## rbdmap

```
# ./rbdmap -h
Usage of ./rbdmap:
  -n, --dry-run            dry run (don't actually map image)
      --id string          Specifies the username (without the 'client.' prefix)
  -i, --image string       Image to map
  -m, --monitor strings    Connect to one or more monitor addresses (192.168.0.1[:6789]). Multiple address are specified comma separated.
      --namespace string   Use a pre-defined image namespace within a pool
  -p, --pool string        Interact with the given pool.
      --read-only          Map the image read-only
      --secret string      Specifies the user authentication secret
      --snap string        Specifies a snapshot name
  -v, --verbose            Enable additional output
```

## rbdunmap

```
# ./rbdunmap -h
Usage of ./rbdunmap:
      --devid int   RBD Device Index (default 0)
  -n, --dry-run     dry run (don't actually unmap device)
      --force       Optional force argument will wait for running requests and then unmap the image
  -v, --verbose     Enable additional output
```