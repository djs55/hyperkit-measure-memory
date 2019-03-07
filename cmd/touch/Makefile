
touch: main.c
	cc -o touch main.c

linuxkit-build-hyperkit:
	linuxkit build touch.yml

linuxkit-run-hyperkit:
	rm -f foo.img
	linuxkit run -disk $(shell pwd)/foo.img,size=2G  touch

linuxkit-build-qemu:
	linuxkit build -format iso-bios touch.yml

linuxkit-run-qemu:
	rm -f foo.qcow2
	linuxkit run qemu -disk $(shell pwd)/foo.qcow2,size=2G -mem 2048 -iso ./touch.iso
