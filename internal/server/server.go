package server

import (
	"log"
	"net"
	"strconv"
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
	addr := ":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", addr)
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
	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!"))
	return
}
