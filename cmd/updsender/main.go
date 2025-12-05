package main

import (
	"net"
	"log"
	"os"
	"bufio"
	"fmt"
)

func main(){
	udp, udpErr := net.ResolveUDPAddr("udp", "localhost:42069")

	if udpErr != nil {

		log.Fatalf("An error occured while resolving UDP: %s\n", udpErr.Error())
		return
	}

	conn, conErr := net.DialUDP("udp", nil, udp)

	if conErr != nil {

		log.Fatalf("An error occured while dialing UDP: %s\n", conErr.Error())
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")

		bytes, _, readErr := reader.ReadLine()
		
		if readErr != nil {
			log.Fatalf("An error occured while reading: %s\n", readErr.Error())
			continue
		}

		_, conWriteErr := conn.Write(bytes)

		if conWriteErr != nil {
			log.Fatalf("An error occured while writing to conn: %s\n", conWriteErr.Error())
			continue
		}
	}
}
