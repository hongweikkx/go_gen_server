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
	startLink(TestMod)
	time.Sleep(3 * time.Second)
	fmt.Println("hello world start")
	var msg com_server.Msg
	msg.Fun = HandleCallHello
	msg.Args = []interface{}{"hello", 1}

	r := com_server.Call(TestMod.ch, msg)
	k := reflect.TypeOf(r)
	fmt.Println("i get the return:", k.Kind())
	time.Sleep(5 * time.Second)
	return
}

func startLink(mod testModule) {
	fmt.Println("test_1 startLink")
	com_server.StartLink(mod, &mod.ch, 1000, nil)
}

func (a testModule) Init(a interface{}){
	fmt.Println("test_1 init", a)
}


func (a testModule) CodeChange(){
	fmt.Println("test_1 code_change")
}

func (a testModule) Terminate(exitReason string){
	fmt.Println("test_1 terminate")
}


func HandleCallHello(a string, b int) string{
	fmt.Println(a, b)
	return "ok"
}
