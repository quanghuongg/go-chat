package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func onMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, _ := reader.ReadString('\n')
		fmt.Print(msg)
	}
}

func ping(conn net.Conn) {

}

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Input name: ")
	nameReader := bufio.NewReader(os.Stdin)
	name, _ := nameReader.ReadString('\n')
	name = name[:len(name)-1]
	fmt.Println("************MSG************")
	go onMessage(conn)
	for {
		msgReader := bufio.NewReader(os.Stdin)
		msg, err := msgReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		msg = msg[:len(msg)-1]
		msg = fmt.Sprintf("%s: %s\n", name, msg)
		conn.Write([]byte(msg))
	}
	conn.Close()

}
