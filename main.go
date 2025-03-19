package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Where we should start reading from
	var offset int64
	buffer := make([]byte, 8)
	for {
		n, err := f.ReadAt(buffer, offset)
		fmt.Printf("read: %s\n", buffer[:n])

		if err == io.EOF {
			break
		}

		offset += int64(n)
	}
}
