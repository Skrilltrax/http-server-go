package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// Handle Request
	defer closeConnection(conn)
	handleRequest(conn)
}

func handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)

	request, err := ParseRequest(reader)
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		return
	}

	var response *Response
	if request.target == "/index.html" {
		response = NewResponse(HTTP11, Success, []string{}, "")
	} else {
		response = NewResponse(HTTP11, NotFound, []string{}, "")
	}

	_, err = conn.Write([]byte(response.String()))
}

func closeConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Println("Error closing connection: ", err.Error())
		os.Exit(1)
	}
}
