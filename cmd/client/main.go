package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/djs55/hyperkit-measure-memory/pkg/mem"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
	"github.com/pkg/errors"
)

func main() {
	var interval int
	flag.IntVar(&interval, "interval", 60, "interval between samples in seconds")
	var results string
	flag.StringVar(&results, "results", "results", "directory to contain results")
	flag.Parse()

	if err := os.Mkdir(results, 0755); err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create results directory: %v", err)
	}
	count := 0
	for {
		path := filepath.Join(results, fmt.Sprintf("%d", count))
		output, err := os.Create(path)
		if err != nil {
			log.Fatalf("Failed to create %s: %v", path, err)
		}
		writeErr := one(output)
		if err := output.Close(); err != nil {
			log.Fatalf("Failed to close %s: %v", path, err)
		}
		time.Sleep(time.Duration(interval) * time.Second)

		if writeErr != nil {
			log.Println(writeErr)
			if err := os.Remove(path); err != nil {
				log.Fatalf("Failed to remove %s: %v", path, err)
			}
			continue
		}
		count++
	}
}

func connect() net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err == nil {
		log.Printf("Failed to connect to server: %v", err)
	}
	return conn
}

func one(output *os.File) error {
	var mi mem.Meminfo

	conn := connect()
	if conn != nil {
		defer conn.Close()

		dec := json.NewDecoder(conn)
		if err := dec.Decode(&mi); err != nil {
			return errors.Wrapf(err, "Unable to decode json")
		}
	}
	ps, err := mem.GetPS()
	if err != nil {
		return errors.Wrapf(err, "Unable to query ps")
	}
	footprint, err := mem.GetFootprint("com.docker.hyperkit")
	if err != nil && err != mem.ErrNoPhysicalFootprint {
		return errors.Wrapf(err, "Unable to query hyperkit footprint")
	}
	firefoxFootprint, err := mem.GetFootprint("firefox")
	if err != nil && err != mem.ErrNoPhysicalFootprint {
		log.Println("Unable to query firefox footprint")
		firefoxFootprint = mem.Footprint(0)
	}
	touchFootprint, err := mem.GetFootprint("touch")
	if err != nil && err != mem.ErrNoPhysicalFootprint {
		log.Println("Unable to query touch footprint")
		touchFootprint = mem.Footprint(0)
	}
	vmstat, err := mem.GetVMStat()
	if err != nil {
		return errors.Wrapf(err, "Unable to query vmstat: %v")
	}
	sample := sample.Sample{
		Time:             time.Now(),
		Meminfo:          mi,
		PS:               ps,
		Footprint:        footprint,
		FirefoxFootprint: firefoxFootprint,
		TouchFootprint:   touchFootprint,
		VMStat:           vmstat,
	}

	enc := json.NewEncoder(output)
	if err := enc.Encode(&sample); err != nil {
		return errors.Wrapf(err, "Unable to write sample")
	}
	return nil
}
