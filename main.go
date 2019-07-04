package main

import (
	"bufio"
	"context"
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

// BareMessage just does fine
type BareMessage struct {
	Message string
}

func main() {

	var (
		inc       = make(chan api.AnnotatedMessage)
		out       = make(chan api.AnnotatedMessage)
		interrupt = make(chan os.Signal)
		err       error
		conn      io.ReadWriteCloser
	)

	conn, err = net.Dial(os.Args[1], "127.0.0.1:12345")
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "some error %v\n", err)
		os.Exit(1)
	}
	err = subscribe(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "some error %v\n", err)
		os.Exit(1)
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)

	go listenAndPrint(conn, inc)
	go listenInput(out)

	go func(cx context.Context) {
	loop:
		for {
			select {
			case <-cx.Done():
				fmt.Println("goodbye!")
				break loop
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
	}(ctx)

	<-interrupt
}

func listenInput(ch chan api.AnnotatedMessage) {
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		pl := BareMessage{
			Message: s.Text(),
		}
		plData, err := json.Marshal(pl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error wrapping message: %v", err)
			continue
		}

		msg := api.AnnotatedMessage{
			Command: "PUB",
			Target:  "topic_1",
			Payload: plData,
		}
		ch <- msg
	}
}

func listenAndPrint(rd io.Reader, ch chan api.AnnotatedMessage) {
	for {
		p := make([]byte, 16384)
		n, err := bufio.NewReader(rd).Read(p)
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stdout, "something happened very bad: %v", err)
		}
		var msg api.AnnotatedMessage
		if err := json.Unmarshal(p[:n], &msg); err != nil {
			fmt.Fprintf(os.Stdout, "something happened: %v\n", err)
			continue
		}
		ch <- msg
	}
}

func subscribe(rd io.Writer) error {
	sub := &api.AnnotatedMessage{
		Command: "SUB",
		Target:  "topic_1",
	}

	b, err := json.Marshal(sub)
	if err != nil {
		return fmt.Errorf("could not send subscribe message: %v", err)
	}
	rd.Write(b)
	return nil
}
