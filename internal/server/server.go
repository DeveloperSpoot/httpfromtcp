package server

import (
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

	writer := response.NewWriter(conn)
	s.handler(writer, req)
}
