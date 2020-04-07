package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type client struct {
	outgoing chan string // an outgoing message channel
	name     string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:

			// Broadcast incoming message to all clients' outgoing message channels
			for cli := range clients {
				cli.outgoing <- msg
			}
		case cli := <-entering:
			clients[cli] = true
			connected := make([]string, 0, len(clients))
			for cli := range clients {
				connected = append(connected, cli.name)
			}
			cli.outgoing <- fmt.Sprintf("Connected Clients: %#v", connected)
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.outgoing)
		default:

		}
	}
}

func handleConn(conn net.Conn) {

	cli := client{
		outgoing: make(chan string), // outgoing client messages
		name: func() string {
			fmt.Fprint(conn, "Enter your name: ")
			n, _ := bufio.NewReader(conn).ReadString('\n')
			return strings.TrimRight(n, "\n")
		}(),
	}
	go clientWriter(conn, cli.outgoing)

	cli.outgoing <- "You are " + cli.name
	messages <- cli.name + " has arrived"
	entering <- cli

	toSend := make(chan string)
	go func() {
		input := bufio.NewScanner(conn)
		for input.Scan() {
			toSend <- cli.name + ": " + input.Text()
		}
	}()

	d := 10 * time.Second
	timer := time.NewTimer(d)
	for {
		select {
		case m := <-toSend:
			messages <- m
			timer.Stop()
			timer.Reset(d)
		case <-timer.C:
			leaving <- cli
			messages <- cli.name + " has left"
			conn.Close()
		}
	}

}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
