package server

import (
	"net"
)

type Server struct {
	state int
}

func (s *Server) Serve(port int) (*Server, error){

}

func s(s *Server) Close() error {

}

func (s *Server) listen(){

}

func (s *Server) handle(conn, net.Conn)