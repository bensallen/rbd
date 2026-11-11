// Package smoke exercises pkg/krbd, pkg/mount, and the rbd CLI against a live
// Ceph cluster and the kernel rbd module. The tests are behind the "smoke"
// build tag, require root, and are normally run inside the VM started by
// test/smoke/run.sh.
package smoke
