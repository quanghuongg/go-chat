package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	connsMap = make(map[net.Conn]string)
	connCh   = make(chan net.Conn)
	closeCh  = make(chan net.Conn)
	msgCh    = make(chan string)
)

func main() {
	server, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server started")

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s conected \n", conn.RemoteAddr().String())

			connsMap[conn] = ""
			connCh <- conn
		}
	}()

	for {
		select {
		case conn := <-connCh:
			go onMessage(conn)
		case msg := <-msgCh:
			fmt.Print(msg)
		case conn := <-closeCh:
			var name string
			if connsMap[conn] == "" {
				name = "Client"
			} else {
				name = connsMap[conn]
			}
			msg := fmt.Sprintf("%s exit \n", name)
			fmt.Printf(msg)
			delete(connsMap, conn)
			publishMessage(nil, msg)
		}

	}
}

func publishMessage(conn net.Conn, msg string) {
	for key, _ := range connsMap {
		if key != conn {
			key.Write([]byte(msg))
		}
	}
}

func onMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msgCh <- msg
		if connsMap[conn] == "" {
			connsMap[conn] = strings.Split(msg, ":")[0]
		}
		publishMessage(conn, msg)
	}
	closeCh <- conn
}
