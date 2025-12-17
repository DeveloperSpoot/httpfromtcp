package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return nil
	}

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

var endOfRead = "\r\n\r\n"

func proxyHandler(w *response.Writer, req *request.Request) {
	proxyTarget := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")

	resp, err := http.Get("https://httpbin.org" + proxyTarget)
	if err != nil {
		handErr := server.HandlerError{}
		handErr.StatusCode = response.StatusError
		handErr.Message = "Internal error while processing."

		log.Fatalf("An error occured while handling proxy: %s\n", err.Error())
		return
	}

	err = w.WriteStatusLine(response.StatusOK)
	if err != nil {
		log.Fatalf("An error occured while writing status line: %s\n", err.Error())
		return
	}

	head := headers.NewHeaders()
	head.SetHeader("connection", "close")
	head.SetHeader("content-type", "text/plain")
	head.SetHeader("transfer-encoding", "chunked")
	head.SetHeader("trailer", "x-content-sha256, x-content-length")

	err = w.WriteHeaders(head)
	if err != nil {
		log.Fatalf("An error occured while writing proxy headers: %s\n", err.Error())
		return
	}

	buff := make([]byte, 1024)
	readIdx := 0

	for {
		if readIdx >= len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		idx, readErr := resp.Body.Read(buff[readIdx:])

		if readErr != nil && readErr != io.EOF {
			log.Fatalf("An error occured while reading proxy body: %s\n", err.Error())
			return
		}
		readIdx += idx
		_, err = w.WriteEncodingChunk(buff[:readIdx])
		if err != nil {
			log.Fatalf("An error occured while writing proxy body: %s\n", err.Error())
			return
		}

		if bytes.Contains(buff, []byte(endOfRead)) || readErr == io.EOF {
			break
		}

	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Fatalf("An error occured while writing encoding ending: %s\n", err.Error())
		return
	}

	sha := fmt.Sprintf("%x", sha256.Sum256(buff[:readIdx]))
	readStr := fmt.Sprintf("%v", readIdx)

	trailers := headers.NewHeaders()
	trailers.SetHeader("x-content-length", readStr)
	trailers.SetHeader("x-content-sha256", sha)

	log.Println(trailers)

	log.Println(fmt.Sprintf("%v", readIdx))
	log.Println(fmt.Sprintf("%x", sha256.Sum256(buff[:readIdx])))

	w.WriteBody([]byte("x-content-length: " + readStr + "\r\nx-content-sha256: " + sha + "\r\n\r\n"))

	return
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
