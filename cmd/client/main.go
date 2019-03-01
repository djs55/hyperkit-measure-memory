package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/djs55/hyperkit-measure-memory/pkg/mem"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
)

func main() {
	if err := os.Mkdir("results", 0755); err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create results directory: %v", err)
	}
	for count := 0; ; count++ {
		path := filepath.Join("results", fmt.Sprintf("%d", count))
		output, err := os.Create(path)
		if err != nil {
			log.Fatalf("Failed to create %s: %v", path, err)
		}
		one(output)
		if err := output.Close(); err != nil {
			log.Fatalf("Failed to close %s: %v", path, err)
		}
		time.Sleep(time.Second)
	}
}

func connect() net.Conn {
	for {
		conn, err := net.Dial("tcp", "127.0.0.1:1234")
		if err == nil {
			return conn
		}
		fmt.Println("Error dialing:", err.Error())
		time.Sleep(time.Second)
	}
}

func one(output *os.File) {
	conn := connect()
	defer conn.Close()

	var mi mem.Meminfo

	dec := json.NewDecoder(conn)
	if err := dec.Decode(&mi); err != nil {
		log.Fatalf("Unable to decode json: %v", err)
	}

	ps, err := mem.GetPS()
	if err != nil {
		log.Fatalf("Unable to query ps: %v", err)
	}
	footprint, err := mem.GetFootprint("com.docker.hyperkit")
	if err != nil && err != mem.ErrNoPhysicalFootprint {
		log.Fatalf("Unable to query hyperkit footprint: %v", err)
	}
	vmstat, err := mem.GetVMStat()
	if err != nil {
		log.Fatalf("Unable to query vmstat: %v", err)
	}
	sample := sample.Sample{
		Time:      time.Now(),
		Meminfo:   mi,
		PS:        ps,
		Footprint: footprint,
		VMStat:    vmstat,
	}

	enc := json.NewEncoder(output)
	if err := enc.Encode(&sample); err != nil {
		log.Fatalf("Unable to write sample: %v", err)
	}
}
