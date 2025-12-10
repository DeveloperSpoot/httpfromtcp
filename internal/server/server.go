package server

import (
	"fmt"
	"log"
	"net"

	"github.com/DeveloperSpoot/httpfromtcp/internal/response"
)

type Server struct {
	state    int
	listener net.Listener
}

const (
	serverClosed = iota
	serverListening
)

func Serve(port int) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	serv := Server{listener: ln, state: serverListening}

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

	response.WriteStatusLine(conn, response.StatusOK)
	header := response.GetDefualtHeaders(0)

	err := response.WriteHeaders(conn, header)
	if err != nil {
		fmt.Printf("An error occured while writing headers: %v\n", err.Error())
	}
}
