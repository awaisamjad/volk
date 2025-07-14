package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"volk/internal/http"
	// "flag"
)

func main() {
	ln, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on port 8000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection ", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	//? build the request into the correct format string
	var requestBuilder strings.Builder

	startLine, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading start line: %v", err)
		return
	}
	requestBuilder.WriteString(startLine)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading header line: %v", err)
			return
		}

		requestBuilder.WriteString(line)

		if line == "\r\n" || line == "\n" {
			break
		}
	}

	req, err := http.NewRequest(requestBuilder.String())
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nBad Request"))
		return
	}

	resp := req.Response()

	// Write response back to client
	_, err = conn.Write([]byte(resp.String()))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
