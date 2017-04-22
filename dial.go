package main

import (
    "fmt"
    "net"
    "log"
    "bufio"
    //"os"
)

func main() {
    server, err := net.Dial("tcp", "northwestern.edu:80")
    if err != nil {
        log.Println(err)
        return
    }

    // listen to user input and send it to the server
	// inputReader := bufio.NewReader(os.Stdin)
	// msg, err := inputReader.ReadString('\n')
	// if err != nil {
	// 	log.Println(err)
	// }
	//send to server
	//conn.Write([]byte(msg))
    fmt.Fprintf(server, "GET / HTTP/1.0\r\n\r\n")

	//listen to msg from server
	incomingReader := bufio.NewReader(server)
	for {
		incoming, err := incomingReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Connection closed by foreign host.\n")
			break
		}
		log.Println(incoming)
	}
}
