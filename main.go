package main
import (
	"fmt"
	"os"
	"strings"
	"io"
	"errors"
)

func main(){
	file, fileErr := os.Open("./messages.txt")
	if fileErr != nil {
		fmt.Errorf("An error occured while opening the file: %w", fileErr)
		return
	}

	for line := range getLinesChannel(file) {	
		fmt.Println("read: " + line)
	}
}

func getLinesChannel(file io.ReadCloser) <-chan string {
	
	go func() {
    		fmt.Println("goroutine started")
		// ...
	}()

	currentLine := ""

	ch := make(chan string)
	
	go func(){
		defer file.Close()
		defer close(ch)
		for {
			read := make([]byte, 8, 8) // Max 8 bytes	
			n, readErr := file.Read(read)

			parts := []string{}
		
			if readErr != nil {
				if errors.Is(readErr, io.EOF) {	
				break
			}
				fmt.Errorf("An error occured while reading: %w", readErr)
				break
			}

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
