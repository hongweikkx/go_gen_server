package protocolmux

import (
	"fmt"
	"testing"
)

type serv struct{}

var globalc chan Msg

func TestPM(t *testing.T) {
	var m Mux
	var s serv
	globalc = make(chan Msg)
	m.Init(s)
	reply := m.Call("hello")
	fmt.Println(reply)
}

func (s serv) ReadTag(m Msg) int64 {
	if m == "hello" {
		return 1
	} else {
		return 2
	}
}

func (s serv) Send(m Msg) {
	globalc <- m
}

func (s serv) Recv() Msg {
	return <-globalc
}
