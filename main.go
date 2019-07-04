package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/isacikgoz/gomq/api"
)

type BareMessage struct {
	Message string
}

func main() {

	var inc chan *api.AnnotatedMessage
	var out chan *api.AnnotatedMessage

	conn, err := net.Dial(os.Args[1], "127.0.0.1:12345")
	defer conn.Close()
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	subscribe(conn)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)

	go func() {
		for {
			select {
			case msg := <-inc:
				var bare BareMessage
				json.Unmarshal(msg.Payload, &bare)
				log.Printf("%s\n", bare.Message)
			case msg := <-out:
				b, err := json.Marshal(msg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error wrapping message: %v", err)
					os.Exit(1)
				}
				conn.Write(b)
			}
		}
	}()

	go listenAndPrint(conn, inc)
	go listenInput(conn, out)

	<-interrupt
}

func listenInput(conn net.Conn, ch chan<- *api.AnnotatedMessage) {
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		text := s.Text()
		if text == "exit" {
			return
		}
		pl := BareMessage{
			Message: s.Text(),
		}
		plData, err := json.Marshal(pl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error wrapping message: %v", err)
			os.Exit(1)
		}
		msg := &api.AnnotatedMessage{
			Command: "PUB",
			Target:  "topic_1",
			Payload: plData,
		}
		ch <- msg
	}
}

func listenAndPrint(conn net.Conn, ch chan<- *api.AnnotatedMessage) {
	for {
		p := make([]byte, 16384)
		n, err := bufio.NewReader(conn).Read(p)
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stdout, "something happened very bad: %v", err)
		}
		var msg api.AnnotatedMessage
		if err := json.Unmarshal(p[:n], &msg); err != nil {
			fmt.Fprintf(os.Stdout, "something happened: %v\n", err)
			fmt.Fprintf(os.Stdout, "message is: %s\n", p)
		}
		ch <- &msg
	}
}

func subscribe(conn net.Conn) {
	sub := &api.AnnotatedMessage{
		Command: "SUB",
		Target:  "topic_1",
	}

	b, err := json.Marshal(sub)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	conn.Write(b)
}
