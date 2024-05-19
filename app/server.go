package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	handlers map[RequestMethod]map[string]func(request Request) *Response
}

func NewServer() *Server {
	handlerMap := make(map[RequestMethod]map[string]func(request Request) *Response)

	for _, method := range GetAllRequestMethods() {
		handlerMap[method] = make(map[string]func(request Request) *Response)
	}

	return &Server{
		handlers: handlerMap,
	}
}

func (s *Server) AddHandler(method RequestMethod, path string, handler func(request Request) *Response) {
	s.handlers[method][path] = handler
}

func (s *Server) Run() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			log.Println("Failed to set reader deadline:", err)
			// do something else, for example create new conn
			return
		}

		// Handle Request
		s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	defer s.closeConnection(conn)
	reader := bufio.NewReader(conn)

	request, err := ParseRequest(reader)
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		return
	}

	var response *Response
	if _, ok := s.handlers[request.method]; !ok {
		fmt.Println("Invalid request method: ", request)
		response = NewResponse(HTTP11, NotFound, map[string]string{}, "")
	}

	var handlerFunc func(request Request) *Response
	paramMap := make(map[string]string)

	for key, value := range s.handlers[request.method] {
		keyArr := strings.Split(key, "/")
		targetArr := strings.Split(request.target, "/")

		if len(keyArr) != len(targetArr) {
			continue
		}

		var isMatch bool
		paramMap = make(map[string]string)

		for i := 0; i < len(keyArr); i++ {
			isPathParam := s.isPathParam(keyArr[i])

			if !isPathParam && keyArr[i] != targetArr[i] {
				isMatch = false
				break
			}

			if !isPathParam && targetArr[i] == keyArr[i] {
				isMatch = true
				continue
			}

			if isPathParam {
				isMatch = true
				paramMap[keyArr[i][1:len(keyArr[i])-1]] = targetArr[i]
			}
		}

		if isMatch {
			handlerFunc = value
			break
		}
	}

	if handlerFunc == nil {
		response = NewResponse(HTTP11, NotFound, map[string]string{}, "")
	} else {
		request.params = paramMap
		response = handlerFunc(*request)
	}

	_, err = conn.Write([]byte(response.String()))
}

func (s *Server) isPathParam(item string) bool {
	return strings.HasPrefix(item, "{") && strings.HasSuffix(item, "}")
}

func (s *Server) closeConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Println("Error closing connection: ", err.Error())
		os.Exit(1)
	}
}

func handleIndexPage(request Request) *Response {
	return NewResponse(HTTP11, Success, map[string]string{}, "")
}

func handleEcho(request Request) *Response {
	strValue := request.params["str"]
	headerMap := make(map[string]string)

	headerMap["Content-Type"] = "text/plain"
	headerMap["Content-Length"] = strconv.Itoa(len(strValue))

	return NewResponse(HTTP11, Success, headerMap, strValue)
}

func handleUserAgent(request Request) *Response {
	strValue := request.headers["User-Agent"]

	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "text/plain"
	headerMap["Content-Length"] = strconv.Itoa(len(strValue))

	return NewResponse(HTTP11, Success, headerMap, strValue)
}

func main() {
	s := NewServer()
	s.AddHandler(GET, "/", handleIndexPage)
	s.AddHandler(GET, "/echo/{str}", handleEcho)
	s.AddHandler(GET, "/user-agent", handleUserAgent)

	s.Run()
}
