package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/DeveloperSpoot/httpfromtcp/internal/request"
	"github.com/DeveloperSpoot/httpfromtcp/internal/response"
)

type Server struct {
	state    int
	listener net.Listener
	handler  Handler
}

type HandlerError struct {
	statusCode response.StatusCode
	message    string
}

func (he HandlerError) Write(w io.Writer) {

	response.WriteStatusLine(w, he.statusCode)
	head := response.GetDefualtHeaders(len(he.message))
	response.WriteHeaders(w, head)
	w.Write([]byte(he.message))
}

const (
	serverClosed = iota
	serverListening
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

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

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	//	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!"))

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("An error occured while reading request: %s\n", err.Error())
	}

	buff := bytes.NewBuffer([]byte{})
	handlerErr := s.handler(buff, req)

	if handlerErr != nil {
		handlerErr.Write(conn)
		return
	}

	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		fmt.Printf("An error occured writing status line: %s\n", err.Error())
		return
	}
	header := response.GetDefualtHeaders(len(buff.Bytes()))

	err = response.WriteHeaders(conn, header)
	if err != nil {
		fmt.Printf("An error occured while writing headers: %s\n", err.Error())
		return
	}

	conn.Write(buff.Bytes())
}
