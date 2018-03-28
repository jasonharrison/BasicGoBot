package main

import (
	"net"
	"testing"
)

func TestInit(t *testing.T) {
	var err error

	bot := &Bot{}
	bot.conn, err = net.Dial("tcp", "irc.freenode.net:6667")
	if err != nil {
		t.Error(err)
	}
}
