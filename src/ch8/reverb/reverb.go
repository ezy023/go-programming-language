package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn, ch chan struct{}) {
	input := bufio.NewScanner(c)
	var wg sync.WaitGroup
	for input.Scan() {
		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			echo(c, text, 1*time.Second)
			ch <- struct{}{}
		}(input.Text())
	}

	defer func() {
		wg.Wait()
		if tc, ok := c.(*net.TCPConn); ok {
			tc.CloseWrite()
		} else {
			c.Close()
		}
	}()

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-ch:
			fmt.Println("Resetting Timer")
		case <-timer.C:
			fmt.Println("TIMEOUT")
			return
		default:
			// do nothing, poll channels
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn, ch)
	}
}
