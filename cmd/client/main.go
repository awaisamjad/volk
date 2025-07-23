package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	requestStr := "GET /products/index.html HTTP/1.1\r\nHost: www.example.com\r\n\r\n"
	fmt.Fprintf(conn, "%s", requestStr)
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
}
