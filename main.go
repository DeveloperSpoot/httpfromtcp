package main
import (
	"fmt"
	"strings"
	"net"
	"io"
	"errors"
	"log"
)

func main(){
	ln, netErr := net.Listen("tcp", ":42069")

	if netErr != nil {
		fmt.Printf("An error occured while attempting to start TCP: %w", netErr)
		return
	}

	defer ln.Close()

	for {
		conn, conErr := ln.Accept()
		if conErr != nil {
			fmt.Printf("An error occured while attempting to accept connection: %w", conErr)
			break
		}

		log.Println("<<<== A Connectioned Has Been Accepted ==>>>")

		for line := range getLinesChannel(conn) {	
			fmt.Println(line)
		}

		log.Println("===>>> Connection Closed <<<===")

	}

}

func getLinesChannel(c net.Conn) <-chan string {
	
	currentLine := ""

	ch := make(chan string)
	
	go func(){
		defer close(ch)

		for {
			read := make([]byte, 8, 8) // Max 8 bytes	
			n, readErr := c.Read(read)

			if readErr != nil {
				if errors.Is(readErr, io.EOF) {return}

				fmt.Println("An Error Occured While Reading From Connection: %w", readErr)
				break
			}
			parts := []string{}

			read = read[:n]
			readString := string(read)

			parts = strings.Split(readString, "\n")

			for i,part := range(parts){
				if i == (len(parts)-1) && len(parts) > 1 {

					ch <- currentLine
					currentLine = ""
				}	 

				currentLine += part
			}
		}
	}()

	return ch

}
