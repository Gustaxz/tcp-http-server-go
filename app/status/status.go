package status

import "fmt"

type HttpStatus struct {
	StatusCode int
	StatusMsg  string
}

var OK = HttpStatus{200, "OK"}
var BadRequest = HttpStatus{400, "Bad Request"}
var InternalServerError = HttpStatus{500, "Internal Server Error"}
var NotFound = HttpStatus{404, "Not Found"}

func FormatStatus(status HttpStatus) []byte {
	return []byte(
		"HTTP/1.1 " + fmt.Sprint(status.StatusCode) + " " + status.StatusMsg + "\r\n\r\n")
}
