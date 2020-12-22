package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	var s Server
	s.Init()
	c1 := make(chan Event)
	err := s.Subscribe(c1)
	if err != nil {
		t.Errorf(err.Error())
	}
	c2 := make(chan Event)
	err = s.Subscribe(c2)
	if err != nil {
		t.Errorf(err.Error())
	}
	s.Publish("hello")
	s.Publish("world")
	a := <-c1
	if a != "hello" {
		t.Errorf("test pubsub error, sub:%+v", a)
	}
	b := <-c2
	if b != "hello" {
		t.Errorf("test pubsub error, sub:%+v", b)
	}
	err = s.Cancel(c1)
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(time.Second)
	c := <-c1
	if c != nil {
		t.Errorf("test pubsub error, sub:%+v", c)
	}
}
