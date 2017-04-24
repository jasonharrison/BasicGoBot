package main

import (
	"fmt"
	"net"
	"bufio"
	"io"
	"log"
	"strings"
)

const (
	nickname string = "testbot"
	channel string = "#botwar"
)

func sendLine(conn net.Conn, s string) {
	fmt.Println("Send: "+s)
	fmt.Fprintf(conn, s+"\n")
}

func loop(conn net.Conn) {
	reader := bufio.NewReader(conn)
	sendLine(conn, "NICK "+nickname)
	sendLine(conn, fmt.Sprintf("USER %[1]s %[1]s %[1]s %[1]s", nickname))
	for {
		in, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Print("Recv: ", in)
		in=strings.TrimSpace(in)
		var ls []string = strings.Split(in, " ")
		if ls[0] == "PING" {
			sendLine(conn, "PONG "+ls[1])
		} else if ls[1] == "376" { // End of MOTD
			sendLine(conn, "JOIN "+channel)
		} else if ls[1] == "PRIVMSG" {
			channel := ls[2]
			msg := strings.Join(ls[3:], " ")[1:]
			fmt.Println(msg)
			if msg == "hello" {
				sendLine(conn, "PRIVMSG "+ channel + " :Hello!")
			}
		}
	}
	reader.Reset(conn)
}

func main() {
	conn, err := net.Dial("tcp", "irc.freenode.net:6667")
	if err != nil {
		panic(err)
	}
	go loop(conn)
	fmt.Scanln() // keep the main thread alive until enter is pressed
}