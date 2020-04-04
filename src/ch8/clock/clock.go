package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Port string

// Doing this the long way for practice
func (p *Port) Set(val string) error {
	*p = Port(val)
	return nil
}

func (p *Port) String() string { return fmt.Sprintf(":%s", *p) }

func main() {
	var port Port = "8000"
	flag.Var(&port, "port", "the port number to run on")
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g. connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g. client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}
