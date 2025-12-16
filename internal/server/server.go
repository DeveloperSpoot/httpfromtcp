package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/DeveloperSpoot/httpfromtcp/internal/headers"
	"github.com/DeveloperSpoot/httpfromtcp/internal/request"
	"github.com/DeveloperSpoot/httpfromtcp/internal/response"
)

type Server struct {
	state    int
	listener net.Listener
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he HandlerError) write(w io.Writer) {

	response.WriteStatusLine(w, he.StatusCode)
	head := response.GetDefualtHeaders(len(he.Message))
	response.WriteHeaders(w, head)
	w.Write([]byte(he.Message))
}

const (
	serverClosed = iota
	serverListening
)

type Handler func(w *response.Writer, req *request.Request) *HandlerError

func Serve(port int, hand Handler) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	serv := Server{listener: ln, state: serverListening, handler: hand}

	go serv.listen()

	return &serv, nil
}

func (s *Server) Close() error {
	s.state = serverClosed

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) listen() {

	for {
		conn, err := s.listener.Accept()
		if err != nil && s.state != serverClosed {
			log.Fatalf("An error occured while accpeting a connection: %s\n", err.Error())
			continue
		}

		if s.state == serverClosed {
			return
		}
		go s.handle(conn)
	}
}

var endOfRead = "\r\n\r\n"

func proxyHandler(w io.Writer, req *request.Request) {
	proxyTarget := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")

	resp, err := http.Get("https://httpbin.org" + proxyTarget)
	if err != nil {
		handErr := HandlerError{}
		handErr.StatusCode = response.StatusError
		handErr.Message = "Internal error while processing."
		handErr.write(w)

		log.Fatalf("An error occured while handling proxy: %s\n", err.Error())
		return
	}

	writer := response.NewWriter(w)

	err = writer.WriteStatusLine(response.StatusOK)
	if err != nil {
		log.Fatalf("An error occured while writing status line: %s\n", err.Error())
		return
	}

	head := headers.NewHeaders()
	head.SetHeader("content-type", "text/plain")
	head.SetHeader("transfer-encoding", "chunked")

	err = writer.WriteHeaders(head)
	if err != nil {
		log.Fatalf("An error occured while writing proxy headers: %s\n", err.Error())
		return
	}

	for {
		buff := make([]byte, 1024, 1024)
		idx, err := resp.Body.Read(buff)

		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatalf("An error occured while reading proxy body: %s\n", err.Error())
			return
		}

		_, err = writer.WriteEncodingChunk(buff[:idx])
		if err != nil {
			log.Fatalf("An error occured while writing proxy body: %s\n", err.Error())
			return
		}

		if bytes.Index(buff, []byte(endOfRead)) != -1 {
			break
		}

	}

	_, err = writer.WriteChunkedBodyDone()
	if err != nil {
		log.Fatalf("An error occured while writing encoding ending: %s\n", err.Error())
		return
	}

	return
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	//	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!"))
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("An error occured while reading request: %s\n", err.Error())

		handErr := HandlerError{}
		handErr.StatusCode = response.StatusBadRequest
		handErr.Message = "Internal error while processing request."

		handErr.write(conn)
		return
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(conn, req)
		return
	}

	writer := response.NewWriter(conn)
	s.handler(writer, req)
}
