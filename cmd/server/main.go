package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/djs55/hyperkit-measure-memory/pkg/mem"
)

func main() {
	log.Println("Starting server")
	l, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on 0.0.0.0:1234")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

const (
	staNano = 0x2000
)

func handleRequest(conn net.Conn) {
	defer conn.Close()
	enc := json.NewEncoder(conn)
	mi, err := mem.GetMeminfo()
	if err != nil {
		log.Printf("Failed reading /proc/meminfo: %v", err)
		return
	}
	if err := enc.Encode(mi); err != nil {
		log.Printf("Failed to write json: %v", err)
		return
	}
}
