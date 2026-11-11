#!/bin/sh
# PID 1 inside the smoke-test VM. RBD_SMOKE_* arrive as kernel cmdline env.
export PATH=/sbin:/bin
mount -t proc proc /proc
mount -t sysfs sysfs /sys
mount -t devtmpfs devtmpfs /dev
mount -t tmpfs tmpfs /run
mount -t tmpfs tmpfs /tmp

fail() {
	echo "$*"
	dmesg | grep -iE 'krb5|ceph|rbd|overlay|virtio|module'
	echo "RBD_SMOKE_RESULT=FAIL"
	echo o >/proc/sysrq-trigger
	sleep 10
}

for m in virtio_net rbd overlay; do
	modprobe "$m" || { fail "failed to load module $m"; exit 1; }
done

ip link set lo up
ip link set eth0 up
ip addr add 10.0.2.15/24 dev eth0
ip route add default via 10.0.2.2

export RBD_SMOKE_BIN=/sbin/rbd
/sbin/smoke.test -test.v -test.timeout 5m >/tmp/test.log 2>&1
rc=$?
while read -r l; do echo "$l"; done </tmp/test.log
if [ $rc -ne 0 ]; then
	fail "smoke tests failed"
elif grep -q -- '--- SKIP' /tmp/test.log; then
	fail "smoke tests skipped"
else
	echo "RBD_SMOKE_RESULT=PASS"
	echo o >/proc/sysrq-trigger
	sleep 10
fi
