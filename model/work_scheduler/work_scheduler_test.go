package workscheduler

import (
	"testing"
	"time"

	"github.com/hongweikkx/go_pattern"
)

func TestWS(t *testing.T) {
	scT, err := NewSchedule(2, 3, go_pattern.NewHandlerFunc(sleepT))
	if err != nil {
		t.Errorf("err:%s", err.Error())
		return
	}
	err = scT.Run()
	if err != nil {
		t.Errorf("err run:%s", err.Error())
	}
}

func sleepT() {
	time.Sleep(time.Second)
}
