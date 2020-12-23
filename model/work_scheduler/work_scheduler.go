// the work scheduler model is used to exec n same task use m goroutine
package workscheduler

import (
	"errors"
	"fmt"

	pattern "github.com/hongweikkx/go_pattern"
)

type Schedule struct {
	numServer int
	numTask   int
	handler   *pattern.HandlerFunc
	srvCh     chan int
}

func NewSchedule(numServer int, numTask int, handle *pattern.HandlerFunc) (*Schedule, error) {
	sc := &Schedule{
		numServer: numServer,
		numTask:   numTask,
		handler:   handle,
	}
	err := sc.checkPre()
	if err != nil {
		return nil, err
	}
	sc.srvCh = make(chan int, sc.numServer)
	for id := 0; id < sc.numServer; id++ {
		sc.srvCh <- id
	}
	return sc, nil
}

func (sc *Schedule) Run() error {
	taskCh := make(chan int, sc.numTask)
	exit := make(chan bool)
	errCh := make(chan error)
	runTasks := func(srvId int) {
		for range taskCh {
			defer func() {
				if err := recover(); err != nil {
					errCh <- errors.New(fmt.Sprintf("[recover]error: %+v", err))
				}
			}()
			sc.handler.Run()
			errCh <- nil
		}
	}

	go func() {
		for {
			select {
			case srv := <-sc.srvCh:
				go runTasks(srv)
			case <-exit:
				return
			}
		}
	}()
	for task := 0; task < sc.numTask; task++ {
		taskCh <- task
	}
	for i := 0; i < sc.numTask; i++ {
		err := <-errCh
		if err != nil {
			close(taskCh)
			exit <- true
			return err
		}
	}
	close(taskCh)
	exit <- true
	return nil
}

func (sc *Schedule) checkPre() error {
	err := sc.handler.Check()
	if err != nil {
		return err
	}
	if sc.numServer <= 0 {
		return errors.New("the server num must bigger than 0")
	}
	if sc.numTask <= 0 {
		return errors.New("the task num must bigger than 0")
	}
	return nil
}
