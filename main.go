package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/isacikgoz/gomq/api"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:12345")
	defer conn.Close()
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := s.Text()
		if text == "exit" {
			return
		}
		pl := api.BareMessage{
			Message: s.Text(),
		}
		plData, err := json.Marshal(pl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error wrapping message: %v", err)
			os.Exit(1)
		}
		msg := &api.AnnotatedMessage{
			Command: "pub",
			Target:  "topic_1",
			Payload: plData,
		}
		b, err := json.Marshal(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error wrapping message: %v", err)
			os.Exit(1)
		}
		conn.Write(b)
	}
}
