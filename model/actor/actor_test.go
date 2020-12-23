package actor

import (
	"testing"
	"time"
)

type testModule Mod

var TestMod testModule

const checkMark = "\u2713"
const ballotX = "\u2717"

func TestComServer(t *testing.T) {
	// ====================== test actor ======================
	TestMod.StartLink()
	t.Log("test_1 startLink", checkMark)
	var msg Msg
	msg.Fun = HandleCallHello
	msg.Args = []interface{}{"hello", 1}
	r := Call(Mod(TestMod), msg)
	t.Log("msg call", checkMark)

	k := ParseRet(r, 0).(string)
	b := ParseRet(r, 1).(int)
	if k != "hello" || b != 1 {
		t.Errorf("i get the return: %v %v", k, ballotX)
	} else {
		t.Log("i get the return:", k, checkMark)
	}

	// cast
	msg.Args = []interface{}{"hello", 1}
	Cast(Mod(TestMod), msg)
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

func HandleCallHello(a string, b int, state interface{}) (string, int, interface{}) {
	return a, b, state.(int) + 1
}
