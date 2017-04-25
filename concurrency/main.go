package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const channel string = "##jason"

type Bot struct {
	conn     net.Conn
	nickname string
}

func (b *Bot) SendLine(s string) {
	fmt.Println(b.nickname + ":send: " + s)
	fmt.Fprintf(b.conn, s+"\n")
}

func (b *Bot) Quit(msg ...string) {
	if len(msg) == 1 {
		b.SendLine("QUIT :" + msg[0])
	} else {
		b.SendLine("QUIT")
	}
}

func (b *Bot) Loop() {
	reader := bufio.NewReader(b.conn)
	b.SendLine("NICK " + b.nickname)
	b.SendLine(fmt.Sprintf("USER %[1]s %[1]s %[1]s %[1]s", b.nickname))
	defer b.Quit()
MAINLOOP:
	for {
		in, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		recv := strings.TrimSpace(in)
		fmt.Println(b.nickname+":recv: ", recv)
		var ls []string = strings.Split(recv, " ")
		if ls[0] == "PING" {
			b.SendLine("PONG " + ls[1])
		} else if ls[1] == "001" { // RPL_WELCOME
			b.SendLine("JOIN " + channel)
		} else if ls[1] == "PRIVMSG" {
			channel := ls[2]
			if strings.Contains(channel, "#") != true { // ignore messages not sent in a channel
				continue MAINLOOP
			}
			msg := strings.Join(ls[3:], " ")[1:]
			fmt.Println(msg)
			if msg == "hello" {
				b.SendLine("PRIVMSG " + channel + " :Hello!")
			} else if msg == b.nickname+": hello" { // nickname: hello
				b.SendLine("PRIVMSG " + channel + " :Hello!")
			}
		}
	}
}

func main() {
	var err error
	var bots []Bot

	bot1 := &Bot{}
	bot1.conn, err = net.Dial("tcp", "irc.freenode.net:6667")
	if err != nil {
		panic(err)
	}
	bot1.nickname = "gobot1"
	bots = append(bots, *bot1)
	go bot1.Loop()

	bot2 := &Bot{}
	bot2.conn, err = net.Dial("tcp", "irc.freenode.net:6667")
	if err != nil {
		panic(err)
	}
	bot2.nickname = "gobot2"
	bots = append(bots, *bot2)
	go bot2.Loop()
	
	fmt.Scanln() // keep the main thread alive until enter is pressed
	for _, b := range bots {
		b.Quit("Keyboard interrupt")
	}
}
