package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	var s Server
	t.Log("init start...")
	s.Init()
	t.Log("init done...")
	c := make(chan Event)
	err := s.Subscribe(c)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log("publish hello...")
	s.Publish("hello")
	t.Log("accept ...")
	a := <-c
	t.Log(a)
	time.Sleep(3 * time.Second)
	err = s.Cancel(c)
	if err != nil {
		fmt.Println(err.Error())
	}
}
