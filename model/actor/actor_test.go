package actor

import (
	"testing"
	"time"

	pattern "github.com/hongweikkx/go_pattern"
)

func TestComServer(t *testing.T) {
	testMod := Start(10)
	now := time.Now().Unix()
	testMod.Cast(pattern.NewHandlerFunc(time.Sleep, 2*time.Second))
	if time.Now().Unix()-now > 1 {
		t.Errorf("cast is not valid")
	}
	sum := 0
	err := testMod.Call(pattern.NewHandlerFunc(add, 1, 2).SetRets(&sum))
	if err != nil {
		t.Errorf("call add error:%s", err.Error())
	}
	if sum != 3 {
		t.Errorf("call add error, sum: %d", sum)
	}
	testMod.Stop("stop")
}

func add(x, y int) int {
	return x + y
}
