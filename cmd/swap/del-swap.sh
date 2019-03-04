#!/bin/sh

for ((i=0;i<8;i++))
do
docker run --privileged --pid=host justincormack/nsenter1 /bin/sh -c "swapoff /var/lib/swap.$i"
docker run --privileged --pid=host justincormack/nsenter1 /bin/sh -c "rm -f /var/lib/swap.$i"
done

