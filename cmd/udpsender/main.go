package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpEndpoint, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	conn, err := net.DialUDP("udp", nil, udpEndpoint)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		in, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %s", err)
		}
		_, err = conn.Write([]byte(in))
		if err != nil {
			log.Fatalf("Error writing to connection: %s", err)
		}
	}
}
