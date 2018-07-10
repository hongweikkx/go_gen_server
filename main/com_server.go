package main

import (
	"reflect"
	"time"
	"fmt"
)

type Msg struct{
	fun interface{}
	args []interface{}
	callRet chan interface{}
}

type GoServer interface{
	startLink(interface{}) ()
	init(...interface{}) ()
	code_change(goServer GoServer)
	terminate(exitReason string)
}

type Module struct{
	callCh chan Msg
	castCh chan Msg
	exitCh chan string
	callRet chan interface{}
}


func StartLink(mod *Module) {
	mod.callCh = make(chan Msg)
	mod.castCh = make(chan Msg, 1000) // the cache can be change
	mod.exitCh = make(chan string)
	mod.callRet = make(chan interface{})
	go doSpawn(*mod)
}

func Call(mod Module, callRet chan interface{}, msg Msg) interface{}{
	msg.callRet = callRet
	mod.callCh <- msg
	select {
	case ret := <- mod.callRet:
	        fmt.Println("call back", ret)
		return ret
	case <- time.After(5 * time.Second):
		panic("call timeout")
	}
}

func Cast(mod Module, msg Msg) {
	mod.castCh <- msg
}

func Stop(mod Module, msg string){
	mod.exitCh <- msg
}

func Reply(mod Module, msg interface{}){
	mod.callRet <- msg
}

func doSpawn(mod Module){
	mod.init()
	loop(mod)
	defer close(mod.callRet)
	defer close(mod.callCh)
	defer close(mod.castCh)
	defer close(mod.exitCh)
}

// 需要增加exit的消息 和 事务的处理
func loop(mod Module){
	select{
	case m1 := <- mod.callCh:
		r := MsgFun(m1)
		Reply(mod, r)
	case m2 := <- mod.castCh:
		MsgFun(m2)
	case exitReason := <- mod.exitCh:
		mod.terminate(exitReason)
		return
	}
	loop(mod)
}

func MsgFun(m Msg) interface{}{
	f := reflect.ValueOf(m.fun)
	in := make([]reflect.Value, len(m.args))
	for k,param := range m.args{
		fmt.Println(param)
		in[k] = reflect.ValueOf(param)
	}
	r := f.Call(in)
	return r
}

