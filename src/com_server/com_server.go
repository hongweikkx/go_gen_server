package com_server

import (
	"reflect"
	"time"
)

type Msg struct{
	Fun interface{}
	Args []interface{}
	CallRet chan interface{}
}

type GoServer interface{
	Init(interface{}) interface{}
	CodeChange(state interface{})
	Terminate(exitReason string, state interface{})
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
	state := mod.Init(opt)
	loop(mod, ch, state)
	defer close(ch.CallRet)
	defer close(ch.CallCh)
	defer close(ch.CastCh)
	defer close(ch.ExitCh)
}

func loop(mod GoServer, ch Chan, state interface{}){
	select{
	case callMsg := <- ch.CallCh:
		r := MsgFun(callMsg, state)
		var ret interface{}
	        ret, state = get_state(r)
		Reply(ch, ret)
	case castMsg := <- ch.CastCh:
		r := MsgFun(castMsg, state)
		_, state = get_state(r)
	case exitReason := <- ch.ExitCh:
		mod.Terminate(exitReason, state)
		return
	}
	loop(mod, ch, state)
}

func MsgFun(m Msg, state interface{}) ([]reflect.Value){
	nl := len(m.Args) +  1
	nargs := make([]interface{}, nl)
	for i:= range m.Args{
		nargs[i] = m.Args[i]
	}
	nargs[nl - 1] = state
	return Apply(m.Fun, nargs)
}

func Apply(f interface{}, args []interface{})([]reflect.Value){
	fun := reflect.ValueOf(f)
	in := make([]reflect.Value, len(args))
	for k,param := range args{
		in[k] = reflect.ValueOf(param)
	}
	r := fun.Call(in)
	return r

}

func get_state(r []reflect.Value)([]reflect.Value, interface{}){
	l := len(r)
	state := r[l - 1]
	return r[:l - 1], state.Interface()
}
