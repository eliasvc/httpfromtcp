package main

import (
	"fmt"
	"log"
	"net"

	"httpfromtcp/internal/request"
)

const bufSize = 8

func main() {

	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error opening port: %q\n", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %s\n", err)
		}
		fmt.Println("Connection accepted")

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for header, value := range r.Headers {
			fmt.Printf("- %s: %s\n", header, value)
		}
	}
}
