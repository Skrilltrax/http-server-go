package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type HttpResponse int

const (
	Success HttpResponse = 200
)

func (h HttpResponse) Value() int {
	return int(h)
}

func (h HttpResponse) String() (string, error) {
	switch h {
	case Success:
		return "Success", nil
	default:
		return "", errors.New("invalid http response")
	}
}

func (h HttpResponse) Reason() (string, error) {
	switch h {
	case Success:
		return "OK", nil
	default:
		return "", errors.New("invalid http response")
	}
}

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

	// Create Status Line
	reason, err := Success.Reason()
	if err != nil {
		fmt.Println("Error getting reason: ", err.Error())
		os.Exit(1)
	}
	statusLine := createStatusLine("HTTP/1.1", strconv.Itoa(Success.Value()), reason)

	response := createResponse(statusLine, []string{}, "")

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}

	err = l.Close()
	if err != nil {
		fmt.Println("Error closing listener: ", err.Error())
		os.Exit(1)
	}
}

func createResponse(statusLine string, headers []string, responseBody string) string {
	sb := strings.Builder{}

	sb.WriteString(statusLine)
	sb.WriteString("\r\n")

	for _, h := range headers {
		sb.WriteString(h)
		sb.WriteString("\r\n")
	}

	sb.WriteString(responseBody)

	return sb.String()
}

func createStatusLine(version string, statusCode string, reason string) string {
	sb := strings.Builder{}

	sb.WriteString(version)
	sb.WriteString(" ")
	sb.WriteString(statusCode)
	sb.WriteString(" ")
	sb.WriteString(reason)

	return sb.String()
}
