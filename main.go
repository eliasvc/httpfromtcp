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
	// Where we should start reading from
	var offset int64
	var currentLine string
	buffer := make([]byte, 8)
	for {
		n, err := f.ReadAt(buffer, offset)
		parts := strings.Split(string(buffer[:n]), "\n")
		if len(parts) == 1 {
			currentLine += parts[0]
		}
		if len(parts) == 2 {
			currentLine += parts[0]
			fmt.Printf("read: %s\n", currentLine)
			currentLine = parts[1]
		}

		if err == io.EOF {
			break
		}
		offset += int64(n)
	}
}
