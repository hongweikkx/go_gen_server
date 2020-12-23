package pattern

import (
	"testing"
	"time"
)

func TestComServer(t *testing.T) {
	testMod := Start(10)
	now := time.Now().Unix()
	testMod.Cast(NewHandlerFunc(time.Sleep, 2*time.Second))
	if time.Now().Unix()-now > 1 {
		t.Errorf("cast is not valid")
	}
	sum := 0
	err := testMod.Call(NewHandlerFunc(add, 1, 2).SetRets(&sum))
	if err != nil {
		t.Errorf("call add error:%s", err.Error())
	}
	if sum != 3 {
		t.Errorf("call add error, sum: %d", sum)
	}
	testMod.Stop("stop")
}
