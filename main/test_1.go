package main

import (
	"fmt"
	"time"
	"reflect"
)


func main(){
	var mod Module
	mod.startLink(&mod)
	time.Sleep(3 * time.Second)
	fmt.Println("hello world start")
	var msg Msg
	msg.fun = HandleCallHello
	msg.args = []interface{}{"hello", 1}
	r := Call(mod, mod.callRet, msg)
	k := reflect.TypeOf(r)
	fmt.Println("i get the return:", k.Kind())
	time.Sleep(5 * time.Second)
	return
}

func (a Module) init(...interface{}){
	fmt.Println("test_1 init")
}

func (a Module) startLink(mod *Module) {
	fmt.Println("test_1 startLink")
	StartLink(mod)
}

func HandleCallHello(a string, b int) string{
	fmt.Println(a, b)
	return "ok"
}

func (a Module) code_change(goServer GoServer){
	fmt.Println("test_1 code_change")
}

func (a Module) terminate(exitReason string){
	fmt.Println("test_1 terminate")
}
