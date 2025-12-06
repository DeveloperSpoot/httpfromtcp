package request

import (
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

func parseRequestLine(r string) (*RequestLine, error) {
	parts := []string{}
	parts = strings.Split(r, "\r\n")

	requestLine := new(RequestLine)

	rl := parts[0]

	rlParts := make([]string, 3, 3)
	rlParts = strings.Split(rl, " ")

	if len(rlParts) < 3 || len(rlParts) > 3 {
		return nil, errors.New("Invalid Request Line: " + rl)
	}

	requestLine.Method = rlParts[0]
	requestLine.RequestTarget = rlParts[1]

	requestLine.HttpVersion = strings.Split(rlParts[2], "/")[1]

	if strings.ToUpper(requestLine.Method) != requestLine.Method {
		return nil, errors.New("Invalid method: " + requestLine.Method)
	}

	if requestLine.Method != "POST" && requestLine.Method != "GET" {
		return nil, errors.New("Unspported method: " + requestLine.Method)
	}

	if requestLine.HttpVersion != "1.1" {
		return nil, errors.New("Invalid or unsupported HTTP version: " + requestLine.HttpVersion)
	}

	return &RequestLine{Method: requestLine.Method, RequestTarget: requestLine.RequestTarget, HttpVersion: requestLine.HttpVersion}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := new(Request)

	bytes, readErr := io.ReadAll(reader)
	if readErr != nil {
		return nil, errors.New("Problem reading from reader: " + readErr.Error())
	}

	rl, rlErr := parseRequestLine(string(bytes))
	if rlErr != nil {
		return nil, rlErr
	}
	request.RequestLine = *rl

	return request, nil
}
