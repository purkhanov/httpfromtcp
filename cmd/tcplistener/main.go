package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Println("read:", line)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}

func getLinesChannel(file io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer file.Close()
		defer close(lines)
		currentLineContents := ""

		for {
			buffer := make([]byte, 8, 8)
			n, err := file.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()

	return lines
}
