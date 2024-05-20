package main

import "errors"

type RequestMethod string

const (
	GET  RequestMethod = "GET"
	POST RequestMethod = "POST"
)

func (RequestMethod RequestMethod) String() string {
	return string(RequestMethod)
}

func ParseRequestMethod(method string) (RequestMethod, error) {
	switch method {
	case "GET":
		return GET, nil
	case "POST":
		return POST, nil
	default:
		return "", errors.New("invalid request method")
	}
}

func GetAllRequestMethods() []RequestMethod {
	return []RequestMethod{GET, POST}
}
