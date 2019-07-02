package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/isacikgoz/gomq/api"
)

type BareMessage struct {
	Message string
}

func main() {

	conn, err := net.Dial(os.Args[1], "127.0.0.1:12345")
	defer conn.Close()
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
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

	go listenAndPrint(conn)

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
		b, err = json.Marshal(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error wrapping message: %v", err)
			os.Exit(1)
		}
		conn.Write(b)
	}
}

func listenAndPrint(conn net.Conn) {
	for {
		p := make([]byte, 16384)
		_, err := bufio.NewReader(conn).Read(p)
		if err != nil {
			fmt.Fprintf(os.Stdout, "something happened very bad: %v", err)
			break
		}
		fmt.Fprintf(os.Stdout, "%s\n", p)
	}
}
