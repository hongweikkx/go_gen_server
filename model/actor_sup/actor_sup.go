//一段废弃的代码： 本意是像用go copy一份erlang 的supervisor的代码。 但是因为没有go pid的改变 所以删除很麻烦。
//如果要做的话，就得需要一个唯一id 和 map机制。 不太好， 至少不够golang
/*
package actorsup

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

	waitGroup sync.WaitGroup // 记录已经启动的group
	children []Child   // 记录已经启动的pid 和他的启动方式 如果是simple_one_for_one 就不记录了
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



func (a supModule) Init(sup interface{}) interface{}{
	nSup := sup.(GoSup)
	var childSpec ChildSpec
	childSpec = nSup.InitSup()
	state := initState(childSpec.stragy, nSup)
	if isSimple(childSpec.stragy.strageType) {
		//init_dynamic(&childSpec.childs)
	} else {
		init_child(childSpec.childs, &state)
	}

	return state
}

func (a supModule) Terminate(exitReason string){
	fmt.Println("test_1 terminate")
}

func (a supModule) StartChild(child Child)bool{
	var msg com_server.Msg
	msg.Fun = SupStartChild
	msg.Args = []interface{}{child}
	return com_server.Call(a.ch, msg) == nil
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

func init_child(childSpec []Child, state *supState) []Child{
	do_start_child(childSpec, state)
}

func do_start_child(children ChildSpec, state *supState){
	l := len(children.childs)
	for i:= 0; i <l; i++{
		com_server.Apply(children.childs[i].startF, children.childs[i].args)
		saveChild(state, children.childs[i])
	}
}

func SupStartChild(child Child, state interface{}) (error, interface{}){
	var err error
	if isSimple(state.(supState).stratgeType){
		err = supStartChildSimple(&state.(supState))
	}else{
		err = supStartChildOther(child, state.(supState))
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

func supStartChildOther(child Child, state *supState)error{
	value := com_server.Apply(child.startF, child.args)
	err := value[0].Interface().(error)
	if err != nil{
		saveChild(state, child)
	}
}

func saveDynamicChild(state *supState) error{
	state.waitGroup.Add(1)
}

func saveChild(state *supState, child Child) error{
	append(state.children, child)
	state.waitGroup.Add(1)
}

func isSimple(s string) bool{
	return s == "simple_one_for_one"
}
*/

// ================ ===================================================================================================
// 另一个废弃的代码  我试图将csp 写成actor 模型， 但是不行
// 1. com_server 一直在接收消息并执行， 但是除非其他人也用我的模型 否则无意义。 比如底层 dial 的实现， 就没有办法融入到其中
// 2. com_sup 我可以用sync.WaitGroup 来收集其他goroute的信息， 但是除了完成信息外 没有其他用。 那为啥还要有这个sup呢？
//            我也可以照着supervisor的代码写， 倒是可以完成对go进程的控制， 但是没有办法移除， 因为没有唯一id， 除非有个注册制。。。
// 可以移除啊， 就像com_sup 需要在这里注册一样， 我可以接受来自所有注册进程的退出信息，并完成重启。原来erlang也是这样的啊。
//package com_sup
//
//import (
//	"com_server"
//	"sync"
//)
//
//type PoolSup struct {
//	work chan com_server.GoServer
//	wg   sync.WaitGroup
//}
//
//func New() *PoolSup {
//	p := PoolSup{
//		work: make(chan com_server.GoServer),
//	}
//	for {
//		w := <-p.work
//		p.wg.Add(1)
//		w.StartLink()
//	}
//	return &p
//}
//
//func (p *PoolSup) Run(g com_server.GoServer) {
//	p.work <- g
//}
//
//func (p *PoolSup) Done() {
//	p.wg.Done()
//}
//
//func (p *PoolSup) Shutdown() {
//	close(p.work)
//	p.wg.Wait()
//}
package actorsup
