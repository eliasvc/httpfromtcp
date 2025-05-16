package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const bufSize = 8

func main() {
	buf := make([]byte, bufSize, bufSize)
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("Error reading the file: %q\n", err)
		}
		if err == io.EOF {
			break
		}

		fmt.Printf("read: %s\n", buf[0:n])
	}

}
