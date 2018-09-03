package com_sup

import (
	"com_server"
	"fmt"
	"sync"
)

var (
	SupMod supModule
)

type supModule struct{
	ch com_server.Chan
}

type GoSup interface {
	StartSup() error
	InitSup() ChildSpec
}

type supState struct{
	sup GoSup
	stratgeType string
	intensity int
	period    int

	waitGroup sync.WaitGroup
	children []Child
	restarts []Child
	dynamicRestarts  int
}

//  返回值
type ChildSpec struct{
	childs []Child
	stragy Stragy
}

//
type Child struct{
	restart string
	shutdown string
	workType string
	startF func()
	args []interface{}
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
func (a supModule) Init(sup interface{}) interface{}{
	nSup := sup.(GoSup)
	var childSpec ChildSpec
	childSpec = nSup.InitSup()
	state := initState(childSpec.stragy, nSup)
	if is_simple(childSpec.stragy.strageType) {
		//init_dynamic(&childSpec.childs)
	} else {
		init_child(childSpec.childs)
	}

	return state
}

func (a supModule) CodeChange(){
	fmt.Println("test_1 code_change")
}

func (a supModule) Terminate(exitReason string){
	fmt.Println("test_1 terminate")
}

// true: 表示成功
func (a supModule) StartChild(g GoSup)bool{
	var msg com_server.Msg
	msg.Fun = SupStartChild
	msg.Args = []interface{}{g}
	return com_server.Call(a.ch, msg) == nil
}

func is_simple(s string) bool{
	return s == "simple_one_for_one"
}

// ========================= init =================================
func initState(s Stragy, g GoSup) supState{
	checkStrage(s)
	state := supState{
		sup:g,
		stratgeType:s.strageType,
		intensity:s.intensity,
		period:s.period,
	}
	return state
}

func checkStrage(s Stragy) {
	b :=
	s.strageType == "simple_one_for_one" ||
	s.strageType == "one_for_one" ||
	s.strageType == "one_for_all" ||
	s.strageType == "rest_for_one"
	if b == false{
		panic("checkStrage error")
	}
}

//func init_dynamic(children *([]Child)){
//}

func init_child(childSpec []Child) []Child{
	do_start_child(childSpec)
}

func do_start_child(children ChildSpec){
	l := len(children.childs)
	for i:= 0; i <l; i++{
		com_server.Apply(children.childs[i].startF, children.childs[i].args)
	}
}

func SupStartChild(g GoSup, state interface{}) (error, interface{}){
	var err error
	if is_simple(state.(supState).stratgeType){
		err = supStartChildSimple(&state.(supState))
	}else{
		err = supStartChildOther(g, state.(supState))
	}
	return err, state
}

func supStartChildSimple(state *supState) error{
	err := state.sup.StartSup()
	if err != nil{
		saveDynamicChild(state)
	}
	return err
}

func supStartChildOther(g GoSup, state *supState)error{
}

func saveDynamicChild(state *supState) error{
}
