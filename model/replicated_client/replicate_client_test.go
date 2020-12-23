package replicatedclient

import (
	"testing"
)

func TestRC(t *testing.T) {
	var client Client
	client.Init([]string{"1", "2"}, func(a string, b Args) Reply {
		t.Log(a, b.(string))
		return "ret"
	})
	client.Call("hello, world")
}
