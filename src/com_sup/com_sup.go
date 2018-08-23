package com_sup

import (
	"com_server"
	"fmt"
)

var (
	SupMod supModule
)

type supModule struct{
	ch com_server.Chan
}

type GoSup interface {
	StartSup()
	InitSup() ChildSpec
}

type ChildSpec struct{
	childs []Child
	stragy Stragy
}
type Child struct{
	restart string
	shutdown string
	workType string
	startF func()
}
type Stragy struct{
	strageType string // simple_one_for_one one_for_one one_for_all
	intensity int
        period int
}


func StartLink(mod supModule,  sup GoSup) {
	fmt.Println("supModule startLink")
	com_server.StartLink(mod, &mod.ch, 1000, sup)
}



// 首先我要知道 这个多个参数该怎么用
// sync
func (a supModule) Init(sup interface{}){
	nSup := sup.(GoSup)
	var childSpec ChildSpec
	childSpec = nSup.InitSup()
	if is_simple(childSpec) {
		init_dynamic(childSpec)
	} else {
		init_child(childSpec)
	}



	fmt.Println("test_1 init")
}

func (a supModule) CodeChange(){
	fmt.Println("test_1 code_change")
}

func (a supModule) Terminate(exitReason string){
	fmt.Println("test_1 terminate")
}

func is_simple(childSpec ChildSpec) bool{
	return childSpec.stragy.strageType == "simple_one_for_one"
}

// ========================= init =================================
func init_dynamic(childSpec ChildSpec) {
	children, err := check_startspec(childSpec)
	if err != nil{
		return "check_error"
	}
}

func init_child(childSpec ChildSpec){
	children, err := check_startspec(childSpec)
	if err != nil {
		return "check_error"
	}
	err = do_start_child(children)
	return err
}

func do_start_child(children )
