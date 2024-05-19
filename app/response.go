package main

import (
	"strconv"
	"strings"
)

type Response struct {
	version Version
	code    ResponseCode
	headers map[string]string
	body    string
}

func NewResponse(version Version, code ResponseCode, headers map[string]string, body string) *Response {
	return &Response{
		version: version,
		code:    code,
		headers: headers,
		body:    body,
	}
}

func (response Response) String() string {
	sb := strings.Builder{}

	// Write Status Line
	statusLine := response.createStatusLine()
	sb.WriteString(statusLine)
	sb.WriteString("\r\n")

	// Write Headers
	for key, value := range response.headers {
		sb.WriteString(key)
		sb.WriteString(": ")
		sb.WriteString(value)
		sb.WriteString("\r\n")
	}
	sb.WriteString("\r\n")

	// Write Response Body
	sb.WriteString(response.body)

	return sb.String()
}

func (response Response) createStatusLine() string {
	sb := strings.Builder{}

	sb.WriteString(response.version.String())
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(response.code.Value()))
	sb.WriteString(" ")
	sb.WriteString(response.code.Reason())

	return sb.String()
}
