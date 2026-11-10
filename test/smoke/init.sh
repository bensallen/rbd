#!/bin/sh
# PID 1 inside the smoke-test VM. RBD_SMOKE_* arrive as kernel cmdline env.
export PATH=/sbin:/bin
mount -t proc proc /proc
mount -t sysfs sysfs /sys
mount -t devtmpfs devtmpfs /dev
mount -t tmpfs tmpfs /run
mount -t tmpfs tmpfs /tmp

for m in virtio_net rbd overlay; do modprobe "$m"; done

ip link set lo up
ip link set eth0 up
ip addr add 10.0.2.15/24 dev eth0
ip route add default via 10.0.2.2

export RBD_SMOKE_BIN=/sbin/rbd
if /sbin/smoke.test -test.v -test.timeout 5m; then
	echo "RBD_SMOKE_RESULT=PASS"
else
	dmesg | grep -iE 'libceph|rbd'
	echo "RBD_SMOKE_RESULT=FAIL"
fi
echo o >/proc/sysrq-trigger
sleep 10
