package main

import (
	"fmt"
	"time"
	"reflect"
	"com_server"
)

var (
	TestMod testModule
)

type testModule struct{
	ch com_server.Chan
}

func main(){
	startLink(&TestMod)

	var msg com_server.Msg
	// call
	msg.Fun = HandleCallHello
	msg.Args = []interface{}{"hello", 1}
	r := com_server.Call(TestMod.ch, msg)
	k := reflect.ValueOf(r)
	fmt.Println("i get the return:", k)
	time.Sleep(3 * time.Second)


	// cast
	msg.Args = []interface{}{"hello", 1}
	com_server.Cast(TestMod.ch, msg)
	time.Sleep(3 * time.Second)

	return
}

func startLink(mod *testModule) {
	fmt.Println("test_1 startLink")
	com_server.StartLink(*mod, &mod.ch, 1000, 1)
}

func (a testModule) Init(i interface{}) interface{}{
	fmt.Println("test_1 init state:", i)
	return i
}


func (a testModule) CodeChange(state interface{}){
	fmt.Println("test_1 code_change")
}

func (a testModule) Terminate(exitReason string, state interface{}){
	fmt.Println("test_1 terminate")
}


func HandleCallHello(a string, b int, state interface{}) (string, interface{}){
	fmt.Println("hello func", a, b, state)
	return "ok", state.(int) + 1
}
