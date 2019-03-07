package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/djs55/hyperkit-measure-memory/pkg/mem"
)

const (
	kib = int64(1024)
	mib = int64(1024) * kib
	gib = int64(1024) * mib

	pageSize = 4096
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
	toAllocate := total + int64(256*mib)
	mem := make([]byte, toAllocate)
	log.Println("Virtual address space is allocated")
	log.Println("Use Control+C to interrupt")
	go func() {
		for {
			for i := 0; i < len(mem); i += pageSize {
				mem[i] = byte(i % 256)
			}
			log.Println("Touched all pages")
			time.Sleep(time.Second)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
