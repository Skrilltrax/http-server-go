package main

import "strings"

type Encoding string

const (
	Gzip    Encoding = "gzip"
	Unknown Encoding = "unknown"
)

func getEncoding(enc string) Encoding {
	enc = strings.ToLower(strings.TrimSpace(enc))

	switch enc {
	case "gzip":
		return Gzip
	default:
		return Unknown
	}
}

func (e Encoding) String() string {
	return string(e)
}
