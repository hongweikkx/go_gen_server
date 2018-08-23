package com_server

import (
	"reflect"
	"time"
	"fmt"
)

type Msg struct{
	Fun interface{}
	Args []interface{}
	CallRet chan interface{}
}

type GoServer interface{
	Init(interface{})
	CodeChange()
	Terminate(exitReason string)
}

type Chan struct{
	CallCh chan Msg
	CastCh chan Msg
	ExitCh chan string
	CallRet chan interface{}
}


// castNum : 异步channel的大小
func StartLink(mod GoServer,ch *Chan, castNum int, opt interface{}) {
	ch.CallCh = make(chan Msg)
	ch.CastCh = make(chan Msg, castNum)
	ch.ExitCh = make(chan string)
	ch.CallRet = make(chan interface{})
	go doSpawn(mod, *ch, opt)
}

func Call(ch Chan, msg Msg) interface{}{
	msg.CallRet = ch.CallRet
	ch.CallCh <- msg
	select {
	case ret := <- ch.CallRet:
	        fmt.Println("call back", ret)
		return ret
	case <- time.After(5 * time.Second):
		panic("call timeout")
	}
}

func Cast(ch Chan, msg Msg) {
	ch.CastCh <- msg
}

func Stop(ch Chan, msg string){
	ch.ExitCh <- msg
}

func Reply(ch Chan, msg interface{}){
	ch.CallRet <- msg
}

func doSpawn(mod GoServer, ch Chan, opt interface{}){
	mod.Init(opt)
	loop(mod, ch)
	defer close(ch.CallRet)
	defer close(ch.CallCh)
	defer close(ch.CastCh)
	defer close(ch.ExitCh)
}

func loop(mod GoServer, ch Chan){
	select{
	case callMsg := <- ch.CallCh:
		r := MsgFun(callMsg)
		Reply(ch, r)
	case castMsg := <- ch.CastCh:
		MsgFun(castMsg)
	case exitReason := <- ch.ExitCh:
		mod.Terminate(exitReason)
		return
	}
	loop(mod, ch)
}

func MsgFun(m Msg) interface{}{
	f := reflect.ValueOf(m.Fun)
	in := make([]reflect.Value, len(m.Args))
	for k,param := range m.Args{
		fmt.Println(param)
		in[k] = reflect.ValueOf(param)
	}
	r := f.Call(in)
	return r
}

