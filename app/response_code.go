package main

type ResponseCode int

const (
	Success             ResponseCode = 200
	Created             ResponseCode = 201
	NotFound            ResponseCode = 404
	InternalServerError ResponseCode = 500
)

func (h ResponseCode) Value() int {
	return int(h)
}

func (h ResponseCode) String() string {
	switch h {
	case Success:
		return "Success"
	case Created:
		return "Created"
	case NotFound:
		return "NotFound"
	case InternalServerError:
		return "InternalServerError"
	default:
		return ""
	}
}

func (h ResponseCode) Reason() string {
	switch h {
	case Success:
		return "OK"
	case Created:
		return "Created"
	case NotFound:
		return "Not Found"
	case InternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}
