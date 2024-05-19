package main

import "errors"

type RequestMethod string

const (
	GET RequestMethod = "GET"
)

func ParseRequestMethod(method string) (RequestMethod, error) {
	switch method {
	case "GET":
		return GET, nil
	default:
		return "", errors.New("invalid request method")
	}
}
