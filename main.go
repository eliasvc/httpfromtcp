package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	lines := make(chan string)
	var currentLine string
	buffer := make([]byte, 8)
	go func() {
		for {
			n, err := f.Read(buffer)
			parts := strings.Split(string(buffer[:n]), "\n")
			// If len == 1 then no \n was encountered
			if len(parts) == 1 {
				currentLine += parts[0]
			}
			// When a \n is encountered, we want to concat whatever we already have in currentLine
			// plus the first part, which forms a complete line. Then we'll want to reset currentLine
			// with the next part. Every part of the slice is actually a complete line, except for the
			// start and the end.
			if len(parts) > 1 {
				currentLine += parts[0]
				for i := 1; i < len(parts); i++ {
					lines <- currentLine
					currentLine = parts[i]
				}
			}

			if err == io.EOF {
				close(lines)
				break
			}
		}
	}()
	return lines
}
