package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type timezone string

func ParseArgs(args []string) map[timezone]string {
	var pairings = make(map[timezone]string)
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		pairings[timezone(parts[0])] = parts[1]
	}
	return pairings
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please supply at least one timezone=hostport pairing")
		os.Exit(1)
	}

	pairings := ParseArgs(os.Args[1:])
	for tz, hp := range pairings {
		conn, err := net.Dial("tcp", hp)
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn, tz)
	}
	for {
		time.Sleep(time.Minute)
	}
}

func handleConn(c net.Conn, tz timezone) {
	defer c.Close()
	lr := bufio.NewReader(c)
	for {
		t, err := lr.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stdout, "%s: %s", tz, t)
	}

}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
