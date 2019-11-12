package actor

import (
	"reflect"
	"testing"
	"time"
)

var (
	TestMod testModule
)

type testModule struct {
	ch Chan
}

const checkMark = "\u2713"
const ballotX = "\u2717"

func TestComServer(t *testing.T) {
	// ====================== test actor ======================
	TestMod.StartLink()
	t.Log("test_1 startLink", checkMark)
	var msg Msg
	msg.Fun = HandleCallHello
	msg.Args = []interface{}{"hello", 1}
	r := Call(TestMod.ch, msg)
	t.Log("msg call", checkMark)

	k := r.([]reflect.Value)[0].Interface().(string)
	if k != "hello" {
		t.Errorf("i get the return: %v %v", k, ballotX)
	} else {
		t.Log("i get the return:", k, checkMark)
	}

	// cast
	msg.Args = []interface{}{"hello", 1}
	Cast(TestMod.ch, msg)
	t.Log("msg cast", checkMark)
	time.Sleep(3 * time.Second)

	return
}

func (mod *testModule) StartLink() {
	StartLink(mod, &mod.ch, 1000, 1)
}

func (mod testModule) Init(i interface{}) interface{} {
	return i
}

func (mod testModule) Terminate(exitReason string, state interface{}) {
}

func HandleCallHello(a string, b int, state interface{}) (string, interface{}) {
	return a, state.(int) + 1
}
