package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/DeveloperSpoot/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK         StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusError      StatusCode = 500
)

func (st *StatusCode) toString() string {
	return strconv.Itoa(int(*st))
}

var statusReason = map[StatusCode]string{
	StatusOK:         "OK",
	StatusBadRequest: "Bad Request",
	StatusError:      "Internal Server Error",
}

var version = "HTTP/1.1"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reason := statusReason[statusCode]

	statusLine := fmt.Sprintf("%v %v %v\r\n", version, statusCode.toString(), reason)

	_, err := w.Write([]byte(statusLine))

	return err
}

func GetDefualtHeaders(contentLen int) headers.Headers {

	head := headers.NewHeaders()

	head.SetHeader("Content-Length", strconv.Itoa(contentLen))
	head.SetHeader("Connection", "close")
	head.SetHeader("Content-Type", "text/plain")

	return head
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for name, value := range headers {
		header := fmt.Sprintf("%v: %v\r\n", name, value)
		_, err := w.Write([]byte(header))

		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))

	return err
}
