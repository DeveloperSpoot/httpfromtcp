package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := new(Request)

	bytes, readErr := io.ReadAll(reader)
	if readErr != nil {
		return nil, errors.New("Problem reading from reader: " + readErr.Error())
	}

	rl, rlErr := parseRequestLine(bytes)
	if rlErr != nil {
		return nil, rlErr
	}
	request.RequestLine = *rl

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, errors.New("No CRLF found in request-line.")
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)

	if err != nil {
		return nil, err
	}

	return requestLine, nil
}

func requestLineFromString(rl string) (*RequestLine, error) {
	rlParts := strings.Split(rl, " ")

	if len(rlParts) < 3 || len(rlParts) > 3 {
		return nil, errors.New("Invalid Request Line: " + rl)
	}

	method := rlParts[0]
	requestTarget := rlParts[1]

	httpParts := strings.Split(rlParts[2], "/")

	protcol := httpParts[0]
	httpVersion := httpParts[1]

	if strings.ToUpper(method) != method {
		return nil, errors.New("Invalid method: " + method)
	}

	if method != "POST" && method != "GET" {
		return nil, errors.New("Unspported method: " + method)
	}

	if strings.Contains(requestTarget, "/") == false {
		return nil, errors.New("Malformed start-line: " + requestTarget)
	}

	if protcol != "HTTP" {
		return nil, errors.New("Invalid Protocol: " + protcol)
	}

	if httpVersion != "1.1" {
		return nil, errors.New("Invalid or unsupported HTTP version: " + httpVersion)
	}

	return &RequestLine{Method: method, RequestTarget: requestTarget, HttpVersion: httpVersion}, nil
}
