package workscheduler

import (
    "testing"
)


func TestWS(t *testing.T) {
    servers := make(chan string, 1)
    servers <- "hello"
    Schedule(servers, 5, func(serv string, task int)bool{
        t.Log(serv, task)
        return true
    })
}

