package main

import "fmt"

type Version string

const (
	HTTP11 Version = "HTTP/1.1"
)

func (v Version) String() string {
	return string(v)
}

func ParseVersion(s string) (Version, error) {
	switch s {
	case "HTTP/1.1":
		return HTTP11, nil
	default:
		return "", fmt.Errorf("unsupported version: %s", s)
	}
}
