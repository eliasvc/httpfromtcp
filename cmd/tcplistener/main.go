package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Println("Connection accepted")
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf(`Request line:
			- Method: %s
			- Target: %s 
			- Version: %s
			`, r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
	}
}
