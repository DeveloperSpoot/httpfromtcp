package main

import (
	"fmt"
	"github.com/DeveloperSpoot/httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	ln, netErr := net.Listen("tcp", ":42069")

	if netErr != nil {
		log.Fatalf("An error occured while attempting to start TCP: %s\n", netErr.Error())
		return
	}

	defer ln.Close()

	for {
		conn, conErr := ln.Accept()
		if conErr != nil {
			log.Fatalf("An error occured while attempting to accept connection: %s\n", conErr.Error())
			break
		}

		log.Println("<<<== A Connectioned Has Been Accepted ==>>>")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("An error occured while getting the request: %s\n", err.Error())
			break
		}

		fmt.Printf("Request line:\n- Method: %v\n- Target: %v\n- Version: %v\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

		log.Println("===>>> Connection Closed <<<===")

	}

}
