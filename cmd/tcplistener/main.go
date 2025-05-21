package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Printf("read: %s\n", line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		buf := make([]byte, bufSize, bufSize)
		var currentLine string
		defer f.Close()
		defer close(out)
		for {
			n, err := f.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatalf("Error reading the file: %q\n", err)
			}
			if err == io.EOF {
				// Dump whatever was read when EOF hit
				if currentLine != "" {
					out <- currentLine
				}
				break
			}
			sBuf := string(buf[0:n])
			parts := strings.Split(sBuf, "\n")
			// The last part will be the one not ending on '\n', so the loop doesn't need
			// to include it
			for i := range len(parts) - 1 {
				currentLine += parts[i]
				out <- currentLine
				currentLine = ""
			}

			// Add the last part, which doesn't include '\n'
			currentLine += parts[len(parts)-1]
		}
	}()
	return out
}
