package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
		fmt.Fprintf(conn, s.Text())
		// p := make([]byte, 2048)
		// _, err = bufio.NewReader(conn).Read(p)
		// if err == nil {
		// 	fmt.Printf("%s\n", p)
		// } else {
		// 	fmt.Printf("Some error %v\n", err)
		// }
	}
}
