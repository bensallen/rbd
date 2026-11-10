#!/bin/bash
# Run the smoke tests in test/smoke against a Ceph demo container, inside a
# QEMU VM booting the host kernel. Requires no root on the host.
#
# Environment:
#   MON_IP             Host IP for Ceph to bind (default: primary IPv4)
#   CEPH_IMAGE         Ceph demo image (default: quay.io/ceph/demo:latest-squid)
#   CONTAINER_RUNTIME  podman or docker (default: podman if present)
#   KERNEL             Kernel image to boot (default: /boot/vmlinuz-$(uname -r))
#   KVER               Module version matching KERNEL (default: $(uname -r))
#   KEEP_CEPH=1        Leave the Ceph container running afterward
set -euo pipefail

REPO=$(cd "$(dirname "$0")/../.." && pwd)
KVER=${KVER:-$(uname -r)}
KERNEL=${KERNEL:-/boot/vmlinuz-$KVER}
CEPH_IMAGE=${CEPH_IMAGE:-quay.io/ceph/demo:latest-squid}
CONTAINER=rbd-smoke-ceph
RT=${CONTAINER_RUNTIME:-$(command -v podman >/dev/null && echo podman || echo docker)}
MON_IP=${MON_IP:-$(ip -4 route get 1.1.1.1 | sed -n 's/.* src \([0-9.]*\).*/\1/p')}
WORK=$(mktemp -d "${TMPDIR:-/tmp}/rbd-smoke.XXXXXX")

log() { echo "==> $*" >&2; }
ceph() { "$RT" exec "$CONTAINER" ceph "$@"; }

cleanup() {
	rm -rf "$WORK"
	[ "${KEEP_CEPH:-}" = 1 ] || "$RT" rm -f "$CONTAINER" >/dev/null 2>&1 || true
}
trap cleanup EXIT

start_ceph() {
	log "starting $CEPH_IMAGE on $MON_IP"
	"$RT" rm -f "$CONTAINER" >/dev/null 2>&1 || true
	"$RT" run -d --name "$CONTAINER" --net=host \
		-e MON_IP="$MON_IP" -e CEPH_PUBLIC_NETWORK="$MON_IP/32" \
		-e DEMO_DAEMONS=mon,mgr,osd "$CEPH_IMAGE" >/dev/null

	for _ in $(seq 60); do
		ceph osd stat 2>/dev/null | grep -q '1 up' && break
		sleep 2
	done
	ceph osd stat | grep -q '1 up' || { "$RT" logs --tail 50 "$CONTAINER"; exit 1; }

	# The demo mon only listens on msgr2; add msgr1 so legacy mode can be tested.
	local mon
	mon=$(ceph mon dump 2>/dev/null | sed -n 's/.* mon\.\(.*\)$/\1/p' | head -1)
	ceph mon set-addrs "$mon" "[v2:$MON_IP:3300,v1:$MON_IP:6789]"

	ceph config set global mon_allow_pool_size_one true
	ceph osd pool create rbd 8 >/dev/null
	ceph osd pool set rbd size 1 --yes-i-really-mean-it >/dev/null
	"$RT" exec "$CONTAINER" rbd pool init rbd
	"$RT" exec "$CONTAINER" rbd create rbd/smoke --size 256M --image-feature layering,exclusive-lock
	ceph auth get-or-create client.smoke mon 'profile rbd' osd 'profile rbd pool=rbd' >/dev/null
	SECRET=$(ceph auth print-key client.smoke)

	for _ in $(seq 30); do
		ceph mon dump 2>/dev/null | grep -q "v1:$MON_IP:6789" && break
		sleep 1
	done
}

build_initramfs() {
	log "building initramfs for $KVER"
	local root=$WORK/root b l ko
	mkdir -p "$root"/{bin,sbin,lib64,dev,proc,sys,run,tmp,newroot}
	ln -s lib64 "$root/lib"

	for b in bash ip mkfs.ext4 kmod mount dmesg grep sleep; do
		b=$(PATH=$PATH:/usr/sbin:/sbin command -v "$b")
		install -D "$b" "$root/sbin/$(basename "$b")"
		ldd "$b" | grep -o '/[^ ]*' | while read -r l; do
			[ -e "$root$l" ] || install -D "$l" "$root$l"
		done
	done
	ln -s bash "$root/sbin/sh"
	ln -s ../sbin/sh "$root/bin/sh"
	ln -s kmod "$root/sbin/modprobe"

	for m in rbd virtio_net overlay; do
		modprobe -S "$KVER" --show-depends "$m"
	done | awk '$1=="insmod"{print $2}' | sort -u | while read -r ko; do
		install -D -m644 "$ko" "$root$ko"
	done
	cp "/lib/modules/$KVER"/modules.{order,builtin,builtin.modinfo} "$root/lib/modules/$KVER/" 2>/dev/null || true
	depmod -b "$root" "$KVER"

	(
		cd "$REPO"
		CGO_ENABLED=0 go build -o "$root/sbin/rbd" ./cmd/rbd
		CGO_ENABLED=0 go test -c -tags smoke -o "$root/sbin/smoke.test" ./test/smoke
	)
	install -m755 "$REPO/test/smoke/init.sh" "$root/init"
	(cd "$root" && find . | cpio -o -H newc --quiet | gzip -1) >"$WORK/initramfs.gz"
}

run_vm() {
	local accel=(-enable-kvm -cpu host)
	if [ ! -w /dev/kvm ]; then
		log "/dev/kvm unavailable, using TCG (slow)"
		accel=(-accel tcg -cpu max)
	fi
	log "booting $KERNEL"
	timeout 600 qemu-system-x86_64 "${accel[@]}" -m 1024 -smp 2 -nographic -no-reboot \
		-kernel "$KERNEL" -initrd "$WORK/initramfs.gz" \
		-netdev user,id=n0 -device virtio-net-pci,netdev=n0 \
		-append "console=ttyS0 quiet panic=-1 RBD_SMOKE_MON=$MON_IP RBD_SMOKE_SECRET=$SECRET" \
		</dev/null | tee "$WORK/console.log"
	grep -q '^RBD_SMOKE_RESULT=PASS' "$WORK/console.log"
}

start_ceph
build_initramfs
run_vm
