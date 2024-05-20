package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type Request struct {
	method  RequestMethod
	target  string
	version Version
	params  map[string]string
	headers map[string]string
	body    []byte
}

func ParseRequest(reader *bufio.Reader) (*Request, error) {
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	method, target, version, err := parseStatusLine(string(line))
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	body := make([]byte, 0)

	// Read headers line by line
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("Error reading line:", err)
			return nil, err
		}

		if len(line) == 0 {
			break
		}

		headerArr, err := parseHeader(string(line))
		if err != nil {
			fmt.Println("Error parsing headers:", err)
			return nil, err
		}

		headers[headerArr[0]] = headerArr[1]
	}

	// Parse body that does not end with \r\n
	sb := strings.Builder{}
	for {
		data, err := reader.ReadByte()
		if err != nil {
			break
		}

		sb.WriteByte(data)
	}

	body = []byte(sb.String())

	return createRequest(method, target, version, headers, body)
}

func createRequest(method RequestMethod, target string, version Version, headers map[string]string, body []byte) (*Request, error) {
	return &Request{
		method:  method,
		target:  target,
		version: version,
		headers: headers,
		body:    body,
	}, nil
}

func parseStatusLine(line string) (RequestMethod, string, Version, error) {
	arr := strings.Split(line, " ")

	method, err := ParseRequestMethod(arr[0])
	if err != nil {
		return "", "", "", err
	}

	target := arr[1]
	version, err := ParseVersion(arr[2])
	if err != nil {
		return "", "", "", err
	}

	return method, target, version, nil
}

func parseHeader(line string) ([]string, error) {
	headerArr := strings.SplitN(line, ":", 2)

	if len(headerArr) != 2 {
		fmt.Println("Invalid header line:", line)
		return nil, errors.New("Invalid header line: " + line)
	}

	headerArr[0] = strings.ToLower(strings.TrimSpace(headerArr[0]))
	headerArr[1] = strings.ToLower(strings.TrimSpace(headerArr[1]))

	return headerArr, nil
}
