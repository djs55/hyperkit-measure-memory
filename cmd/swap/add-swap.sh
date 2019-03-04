#!/bin/sh

for ((i=0;i<8;i++))
do
sleep 30
echo Adding $i
docker run --privileged --pid=host justincormack/nsenter1 /bin/sh -c "dd if=/dev/zero of=/var/lib/swap.$i bs=1024 count=1048576"
docker run --privileged --pid=host justincormack/nsenter1 /bin/sh -c "mkswap /var/lib/swap.$i"
docker run --privileged --pid=host justincormack/nsenter1 /bin/sh -c "swapon /var/lib/swap.$i"
done
