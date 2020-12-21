package util

import (
	"testing"
	"time"
)

func TestParaller(t *testing.T) {
	p1 := NewParaller()
	addRet := 0
	var convertA string
	var convertB int
	var convertC bool
	delRet := 0
	p1.Add(add, 1, 2).SetRets(&addRet)
	p1.Add(convert, "hello", 88, false).SetRets(&convertA, &convertB, &convertC)
	p1.Add(del, 5, 3).SetRets(&delRet)
	err := p1.Run()
	if err != nil || addRet != 3 || delRet != 2 || convertA != "hello" || convertB != 88 || convertC != false {
		t.Error("p1 test err")
	}

	p2 := NewParaller()
	t2Start := time.Now()
	// default 2 * second
	p2.Add(sleep)
	err = p2.Run()
	t2End := time.Now()
	if err != nil || t2End.Unix()-t2Start.Unix() <= 1 {
		t.Error("p2 test err")
	}

	p3 := NewParaller()
	p3.Add(panicX)
	err = p3.Run()
	if err == nil {
		t.Error("p3 test err")
	}
}

func add(x, y int) int {
	return x + y
}

func del(x, y int) int {
	return x - y
}

func convert(a, b, c interface{}) (string, int, bool) {
	aa := a.(string)
	bb := b.(int)
	cc := c.(bool)
	return aa, bb, cc
}

func sleep() {
	time.Sleep(2 * time.Second)
}

func panicX() {
	panic("hello")
}
