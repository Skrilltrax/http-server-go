package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"net"
	"os"
	"strconv"
	"strings"
)

type Context struct {
	directory string
}

type Server struct {
	handlers map[RequestMethod]map[string]func(request Request, ctx Context) *Response
	ctx      Context
}

func NewServer(directory string) *Server {
	handlerMap := make(map[RequestMethod]map[string]func(request Request, ctx Context) *Response)

	for _, method := range GetAllRequestMethods() {
		handlerMap[method] = make(map[string]func(request Request, ctx Context) *Response)
	}

	ctx := Context{directory: directory}

	return &Server{
		handlers: handlerMap,
		ctx:      ctx,
	}
}

func (s *Server) AddHandler(method RequestMethod, path string, handler func(request Request, ctx Context) *Response) {
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

		// Handle Request
		go s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	defer s.closeConnection(conn)

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	reader := bufio.NewReader(bytes.NewReader(buffer))

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

	var handlerFunc func(request Request, ctx Context) *Response
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
		response = handlerFunc(*request, s.ctx)
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

func handleIndexPage(Request, Context) *Response {
	return NewResponse(HTTP11, Success, map[string]string{}, "")
}

func handleEcho(request Request, _ Context) *Response {
	strValue := request.params["str"]
	headerMap := make(map[string]string)

	headerMap["Content-Type"] = "text/plain"
	headerMap["Content-Length"] = strconv.Itoa(len(strValue))

	return NewResponse(HTTP11, Success, headerMap, strValue)
}

func handleUserAgent(request Request, _ Context) *Response {
	strValue := request.headers["User-Agent"]

	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "text/plain"
	headerMap["Content-Length"] = strconv.Itoa(len(strValue))

	return NewResponse(HTTP11, Success, headerMap, strValue)
}

func handleGetFiles(request Request, ctx Context) *Response {
	fileName := request.params["fileName"]

	fileInfo, err := os.Stat(ctx.directory + "/" + fileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("The file does not exist.")
		} else {
			fmt.Println(err)
		}

		return NewResponse(HTTP11, NotFound, map[string]string{}, "")
	}

	file, err := os.Open(ctx.directory + "/" + fileName)
	if err != nil {
		fmt.Println("Error opening file: " + err.Error())
		return NewResponse(HTTP11, NotFound, map[string]string{}, "")
	}

	fileByteArr := make([]byte, fileInfo.Size())
	_, err = bufio.NewReader(file).Read(fileByteArr)
	if err != nil {
		fmt.Println("Error reading file: " + err.Error())
		return NewResponse(HTTP11, NotFound, map[string]string{}, "")
	}

	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "application/octet-stream"
	headerMap["Content-Length"] = strconv.Itoa(int(fileInfo.Size()))

	return NewResponse(HTTP11, Success, headerMap, string(fileByteArr))
}

func handlePostFiles(request Request, ctx Context) *Response {
	fileName := request.params["fileName"]
	filePath := ctx.directory + "/" + fileName

	err := os.MkdirAll(ctx.directory, 0777)
	if err != nil {
		fmt.Println("Error creating directory: " + err.Error())
		return NewResponse(HTTP11, InternalServerError, map[string]string{}, "")
	}

	byteArr := request.body
	byteArr = bytes.Trim(byteArr, "\x00")

	err = os.WriteFile(filePath, byteArr, fs.ModePerm)
	if err != nil {
		fmt.Println("Error writing to file: " + err.Error())
		return NewResponse(HTTP11, InternalServerError, map[string]string{}, "")
	}

	return NewResponse(HTTP11, Created, map[string]string{}, "")
}

func main() {
	var directory string
	if len(os.Args) > 1 {
		directory = os.Args[2]
	} else {
		directory, _ = os.Getwd()
	}

	s := NewServer(directory)
	s.AddHandler(GET, "/", handleIndexPage)
	s.AddHandler(GET, "/echo/{str}", handleEcho)
	s.AddHandler(GET, "/user-agent", handleUserAgent)
	s.AddHandler(GET, "/files/{fileName}", handleGetFiles)
	s.AddHandler(POST, "/files/{fileName}", handlePostFiles)

	s.Run()
}
