package pubsub

// copy from http://nil.csail.mit.edu/6.824/2018/notes/gopattern.pdf

import (
	"errors"
	"fmt"
)

type Event interface{}

type PubSub interface {
	// Publish publishes the event e to all current subscriptions.
	Publish(e Event)
	// Subscribe registers c to receive future events. All subscribers receive events in the same order,
	// and that order respects program order: if Publish(e1) happens before Publish(e2), subscribers receive e1 before e2.
	Subscribe(c chan<- Event)
	// Cancel cancels the prior subscription of channel c.  After any pending already-published events have been sent on c,
	// the server will signal that the subscription is cancelled by closing c.
	Cancel(c chan<- Event)
}

type Server struct {
	publish   chan Event
	subscribe chan subReq
	cancel    chan subReq
}

type subReq struct {
	c  chan<- Event
	ok chan bool
}

func (s *Server) Init() {
	s.publish = make(chan Event)
	s.subscribe = make(chan subReq)
	s.cancel = make(chan subReq)
	go s.loop()
}

func (s *Server) Publish(e Event) {
	s.publish <- e
}

func (s *Server) Subscribe(c chan<- Event) error {
	r := subReq{c: c, ok: make(chan bool)}
	s.subscribe <- r
	if !<-r.ok {
		return errors.New("pubsub: already subscribed")
	}
	return nil
}

func (s *Server) Cancel(c chan<- Event) error {
	r := subReq{c: c, ok: make(chan bool)}
	s.cancel <- r
	if !<-r.ok {
		return errors.New("pubsub: not subscribed")
	}
	return nil
}

func (s *Server) loop() {
	sub := make(map[chan<- Event]chan<- Event)
	for {
		select {
		case e := <-s.publish:
			for _, h := range sub {
				h <- e
			}
		case r := <-s.subscribe:
			if sub[r.c] != nil {
				r.ok <- false
				break
			}
			h := make(chan Event)
			fmt.Println("aaa")
			go helper(h, r.c)
			sub[r.c] = h
			r.ok <- true
		case r := <-s.cancel:
			if sub[r.c] == nil {
				r.ok <- false
				break
			}
			close(sub[r.c])
			delete(sub, r.c)
			r.ok <- true
		}
	}
}

func helper(in <-chan Event, out chan<- Event) {
	q := []Event{}
	for in != nil {
		// http://nil.csail.mit.edu/6.824/2018/notes/gopattern.pdf is in != nil && len(q) > 0 是错误的
		// Decide whether and what to send
		var sendOut chan<- Event
		var next Event
		if len(q) > 0 {
			sendOut = out
			next = q[0]
		}
		select {
		case e, ok := <-in:
			if !ok {
				in = nil
				break
			}
			q = append(q, e)
		case sendOut <- next:
			q = q[1:]
		}
	}
	close(out)
}
