package main

type ResponseCode int

const (
	Success  ResponseCode = 200
	NotFound ResponseCode = 404
)

func (h ResponseCode) Value() int {
	return int(h)
}

func (h ResponseCode) String() string {
	switch h {
	case Success:
		return "Success"
	case NotFound:
		return "NotFound"
	default:
		return ""
	}
}

func (h ResponseCode) Reason() string {
	switch h {
	case Success:
		return "OK"
	case NotFound:
		return "Not Found"
	default:
		return ""
	}
}
