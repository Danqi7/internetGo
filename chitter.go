package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"bufio"
	"strings"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run chitter.go [-c] [port]")
		return
	}
	//connTable maps connection to its client ID
	connTable := make(map[net.Conn]int)
	connCh := make(chan net.Conn)
	deadConnCh := make(chan net.Conn)
	msgCh := make(chan string)
	whoCh := make(chan net.Conn)
	privateCh := make(chan string)
	id := 1

	//start as client
	if os.Args[1] == "-c" {
		client(os.Args)
		return
	}

	//start as server
	port := fmt.Sprintf(":%v", os.Args[1])
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("A error has occured when listening: %v", err.Error())
		os.Exit(1)
	}
	log.Printf("Listening at port: %v", port)

	//listen to incoming connections and pass connections to channel
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			connCh <- conn
		}
	}()

	for {
		select {
		case conn := <- connCh:
			connTable[conn] = id
			go handleClient(conn, id, msgCh, deadConnCh, whoCh, privateCh)
			id++
		case msg := <- msgCh:
			for conn, _ := range connTable {
				conn.Write([]byte(msg))
			}
		case conn := <- whoCh:
			selfID, ok := connTable[conn]
			if !ok {
				break
			}
			conn.Write([]byte(fmt.Sprintf("chitter: %v\n", selfID)))
		case msg:= <- privateCh:
			i := strings.Index(msg, ":")
			from, err := strconv.Atoi(msg[:i])
			if err != nil {
				break
			}
			msg = msg[i+1:]
			i = strings.Index(msg, ":")
			sent_to, err := strconv.Atoi(msg[:i])
			if err != nil {
				break
			}
			msg = msg[i+1:]

			//send to client with the id
			for conn, connID := range connTable {
				if connID == sent_to {
					conn.Write([]byte(fmt.Sprintf("%v: %v", from, msg)))
				}
			}
		case deadConn := <-deadConnCh:
			log.Printf("Client %v is dead", connTable[deadConn])
			delete(connTable, deadConn)
		}
	}

	// close connection when server stops
	for conn, _ := range connTable {
		conn.Close()
	}
}

func handleClient(conn net.Conn, id int, msgCh chan string, deadConnCh chan net.Conn, whoCh chan net.Conn, privateCh chan string) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		i := strings.Index(msg, ":")
		if i == -1 {
			//default command send to all
			start := 0
			for index, c := range msg {
				if string(c) != " " {
					start = index
					break
				}
			}
			msg = msg[start:]
			msgCh <- fmt.Sprintf("%v: %v", id, msg)
			continue
		}
		cmd := strings.Replace(msg[:i], " ", "", -1) //strip all whitespace
		start := i+1      //strip all whitespace before the first appearance of char
		msg = msg[i+1:]
		for index, c := range msg {
			if string(c) != " " {
				start = index
				break
			}
		}
		msg = msg[start:]
		switch cmd {
		case "all":
			msgCh <- fmt.Sprintf("%v: %v", id, msg)
		case "whoami":
			whoCh <- conn
		default:
			_, err := strconv.Atoi(cmd)
			//ignore invalid command
			if err != nil {
				break
			}
			//if 1 send to 3, it would be "1:3:...."
			privateCh <- fmt.Sprintf("%v:%v:%v", id, cmd, msg)
		}

	}
	deadConnCh <- conn
	conn.Close()
}

func client(Args []string) {
	port := fmt.Sprintf(":%v", os.Args[2])
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Printf("Cannot connect to port: %v\n", port)
		os.Exit(1)
	}
	fmt.Printf("Connected to port: %v\n", port)

	// listen to user input and send it to the server
	go func(conn net.Conn) {
		inputReader := bufio.NewReader(os.Stdin)
		for {
			msg, err := inputReader.ReadString('\n')
			if err != nil {
				continue
			}
			//send to server
			conn.Write([]byte(msg))
		}
	}(conn)

	//listen to msg from server
	incomingReader := bufio.NewReader(conn)
	for {
		incoming, err := incomingReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Connection closed by foreign host.\n")
			break
		}
		fmt.Print(incoming)
	}
}
