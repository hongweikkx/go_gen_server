package actor

import (
	"reflect"
	"time"
)

type Msg struct {
	Fun     interface{}
	Args    []interface{}
	CallRet chan interface{}
}

type GoServer interface {
	StartLink()
	Init(interface{}) interface{}
	Terminate(exitReason string, state interface{})
}

type Mod struct {
	ch Chan
}

type Chan struct {
	CallCh  chan Msg
	CastCh  chan Msg
	ExitCh  chan string
	CallRet chan interface{}
}

// castNum : 异步channel的大小
func StartLink(mod GoServer, ch *Chan, castNum int, opt interface{}) {
	ch.CallCh = make(chan Msg, 1)
	ch.CastCh = make(chan Msg, castNum)
	ch.ExitCh = make(chan string)
	ch.CallRet = make(chan interface{})
	go doSpawn(mod, *ch, opt)
}

func Call(mod Mod, msg Msg) interface{} {
	ch := mod.ch
	msg.CallRet = ch.CallRet
	ch.CallCh <- msg
	select {
	case ret := <-ch.CallRet:
		return ret
	case <-time.After(5 * time.Second):
		panic("call timeout")
	}
}

func Cast(mod Mod, msg Msg) {
	mod.ch.CastCh <- msg
}

func Stop(ch Chan, msg string) {
	ch.ExitCh <- msg
}

func Reply(ch Chan, msg interface{}) {
	ch.CallRet <- msg
}

func doSpawn(mod GoServer, ch Chan, opt interface{}) {
	defer close(ch.CallRet)
	defer close(ch.CallCh)
	defer close(ch.CastCh)
	defer close(ch.ExitCh)
	state := mod.Init(opt)
	loop(mod, ch, state)
}

func loop(mod GoServer, ch Chan, state interface{}) {
	select {
	case callMsg := <-ch.CallCh:
		r := msgFun(callMsg, state)
		var ret interface{}
		ret, state = getState(r)
		Reply(ch, ret)
	case castMsg := <-ch.CastCh:
		r := msgFun(castMsg, state)
		_, state = getState(r)
	case exitReason := <-ch.ExitCh:
		mod.Terminate(exitReason, state)
		return
	}
	loop(mod, ch, state)
}

func msgFun(m Msg, state interface{}) []reflect.Value {
	nl := len(m.Args) + 1
	nargs := make([]interface{}, nl)
	for i := range m.Args {
		nargs[i] = m.Args[i]
	}
	nargs[nl-1] = state
	return apply(m.Fun, nargs)
}

func apply(f interface{}, args []interface{}) []reflect.Value {
	fun := reflect.ValueOf(f)
	in := make([]reflect.Value, len(args))
	for k, param := range args {
		in[k] = reflect.ValueOf(param)
	}
	r := fun.Call(in)
	return r

}

func getState(r []reflect.Value) ([]reflect.Value, interface{}) {
	l := len(r)
	state := r[l-1]
	return r[:l-1], state.Interface()
}

// ---------------------------- util func -------------------------
func ParseRet(ret interface{}, i int) interface{} {
	return ret.([]reflect.Value)[i].Interface()
}
