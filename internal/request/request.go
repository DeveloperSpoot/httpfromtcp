package request

import (
	"bytes"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/DeveloperSpoot/httpfromtcp/internal/headers"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	ParserState int
	Headers     headers.Headers
	Body []byte
}

const (
	requestInialized int = iota
	requestDone
	requestParsingHeaders
	requestCheckingBody
	requestParsingBody
)

const bufferSize int = 8
const crlf string = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	request := new(Request)

	request.ParserState = requestInialized
	request.Headers = make(headers.Headers)

	for request.ParserState != requestDone {

		if readToIndex >= len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		bytesRead, readErr := reader.Read(buff[readToIndex:])

		if bytesRead == 0 && errors.Is(readErr, io.EOF) && readToIndex == 0{
			if(request.ParserState == requestParsingBody){
				return nil, errors.New("Body content is less than reported content-length")
			}

			request.ParserState = requestDone
			break
		}
		
		readToIndex += bytesRead
		
		parsed, err := request.parse(buff[:readToIndex])
		if err != nil {
			return nil, errors.New("Error occured while parsing: " + err.Error())
		}
		copy(buff, buff[parsed:])

		readToIndex -= parsed
	}

	return request, nil
}

func (request *Request) parse(data []byte) (int, error) {
	if request.ParserState == requestDone {
		return 0, errors.New("Attetmped to parse request that is done.")
	}

	switch request.ParserState {
	case requestDone:
		return 0, errors.New("Attempted to parse request that is done.")

	case requestInialized:
		idx, requestLine, err := parseRequestLine(data)

		if err == nil && idx == 0 && requestLine == nil {

			return 0, nil
		}

		if err != nil {
			return 0, err
		}

		request.RequestLine = *requestLine
		request.ParserState = requestParsingHeaders
		return idx, nil

	case requestParsingHeaders:
		idx, done, err := request.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			request.ParserState = requestCheckingBody
			return idx, nil
		}
		return idx, nil

	case requestCheckingBody:
			if request.Headers["content-length"] == "" {
				request.ParserState = requestDone
			return 0, nil
		}

		request.ParserState = requestParsingBody
		return 0, nil

	case requestParsingBody:

		leng, err := strconv.ParseInt(request.Headers["content-length"], 0, 0)

		if err != nil {
			return 0, err
		}
		request.Body = append(request.Body, data...)

		if len(request.Body) > int(leng){
			return 0, errors.New("Content-Length does not match the body length.")
		}

		if len(request.Body) < int(leng){
			return len(data), nil
		}

		request.ParserState = requestDone
		log.Println("<><><><>")
		log.Println(request.Body)
		return 0, nil

	default:
		return 0, errors.New("Attempted To Parse During Unknown Parser State")

	}
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))

	// -1 indicates that bytes passed do not contain all of the request-line,
	// so it returns and trys again on the next call. Once it has all of the request-line, it'll parse.
	if idx == -1 {
		return 0, nil, nil
	}

	requestLineText := string(data[:idx])

	requestLine, err := requestLineFromString(requestLineText)

	if err != nil {
		return 0, nil, err
	}

	return len(data[:idx]) + 1, requestLine, nil
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
