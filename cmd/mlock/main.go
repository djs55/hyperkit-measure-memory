package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/djs55/hyperkit-measure-memory/pkg/mem"
)

// Needs to be run with --privileged

const (
	kib = int64(1024)
	mib = int64(1024) * kib
	gib = int64(1024) * mib
)

func main() {
	mi, err := mem.GetMeminfo()
	if err != nil {
		log.Fatalf("Failed to read /proc/meminfo: %v", err)
	}
	total, ok := mi["MemTotal"]
	if !ok {
		log.Fatal("Failed to find key MemTotal in /proc/meminfo")
	}
	log.Printf("System has %d MiB of RAM according to /proc/meminfo", total/mib)
	log.Println("Attempting to allocate everything except the last 512 MiB")
	toAllocate := total - int64(512)*mib
	mem := make([]byte, toAllocate)
	log.Println("Virtual address space is allocated")
	if err := syscall.Mlock(mem); err != nil {
		log.Println("Unable to lock memory. Make sure you run with --privileged")
		log.Fatal(err)
	}
	log.Println("Memory is locked to RAM")
	log.Println("Use Control+C to interrupt")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
