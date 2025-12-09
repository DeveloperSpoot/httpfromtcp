package server

import (
	"net"
)

type Server struct {
	state int
	listener net.Listener
}

const (
	serverClosed = iota
	serverListening
)

func Serve(port int) (*Server, error){
	addr := ":"+string(port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{listener: ln, state: 1}, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()

	return err
}

func (s *Server) listen(){

}

func (s *Server) handle(conn, net.Conn){
	
}