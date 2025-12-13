package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DeveloperSpoot/httpfromtcp/internal/headers"
	"github.com/DeveloperSpoot/httpfromtcp/internal/request"
	"github.com/DeveloperSpoot/httpfromtcp/internal/response"
	"github.com/DeveloperSpoot/httpfromtcp/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) *server.HandlerError {
	head := headers.NewHeaders()
	head.GetDefualtHeaders(0)
	head.SetHeader("content-type", "text/html")

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		body := `
		<html>
 		 <head>
		    <title>400</title>
		  </head>
 		 <body>
		    <h1>Bad Request</h1>
		    <p>Your request honestly kinda sucked.</p>
		 </body>
		</html>
		`
		w.WriteStatusLine(response.StatusBadRequest)

		head.SetHeader("content-length", fmt.Sprintf("%v", len(body)))
		w.WriteHeaders(head)

		w.WriteBody([]byte(body))

	case "/myproblem":
		body := `
		<html>
 		 <head>
		    <title>500 Internal Server Error</title>
		  </head>
 		 <body>
		    <h1>Internal Server Error</h1>
		    <p>Okay, you know what? This one is on me.</p>
		 </body>
		</html>
		`
		w.WriteStatusLine(response.StatusError)

		head.SetHeader("content-length", fmt.Sprintf("%v", len(body)))
		w.WriteHeaders(head)

		w.WriteBody([]byte(body))

	default:
		body := `
		<html>
 		 <head>
		    <title>200 OK</title>
		  </head>
 		 <body>
		    <h1>Success!</h1>
		    <p>Your request was an absolute banger.</p>
		 </body>
		</html>
		`
		w.WriteStatusLine(response.StatusOK)

		head.SetHeader("content-length", fmt.Sprintf("%v", len(body)))
		w.WriteHeaders(head)

		w.WriteBody([]byte(body))

	}

	return nil
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
