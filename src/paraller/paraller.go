// 用于并行执行不同的函数
package util

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type HandlerFunc struct {
	F    interface{}
	Args []interface{}
	Ret  []interface{}
}

type Paraller struct {
	wg       sync.WaitGroup
	handlers []*HandlerFunc
	ErrChan  chan error
}

func NewParaller() *Paraller {
	return &Paraller{}
}

func (p *Paraller) Add(f interface{}, args ...interface{}) *HandlerFunc {
	handle := HandlerFunc{
		F:    f,
		Args: args,
		Ret:  []interface{}{},
	}
	p.handlers = append(p.handlers, &handle)
	return &handle
}

func (p *Paraller) Run() error {
	if err := p.checkPre(); err != nil {
		return err
	}
	if hl := len(p.handlers); hl > 0 {
		p.ErrChan = make(chan error, hl)
		p.wg.Add(len(p.handlers))
		for k := range p.handlers {
			go func(f HandlerFunc, ch chan error) {
				defer p.wg.Done()
				defer func() {
					if err := recover(); err != nil {
						ch <- errors.New(fmt.Sprintf("[recover]error: %+v", err))
					}
				}()
				f.run()
			}(*p.handlers[k], p.ErrChan)
		}
		p.wg.Wait()
	}
	return p.errorRet()
}

func (p *Paraller) checkPre() error {
	for _, handler := range p.handlers {
		for _, ret := range handler.Ret {
			if reflect.ValueOf(ret).Type().Kind() != reflect.Ptr {
				return errors.New("the ret is not ptr")
			}
		}
		fT := reflect.ValueOf(handler.F).Type()
		if fT.Kind() != reflect.Func {
			return errors.New("the func is not function")
		}
		if fT.NumIn() != len(handler.Args) {
			return errors.New("the func args is not valid")
		}
	}
	return nil
}

func (p *Paraller) errorRet() error {
	if len(p.ErrChan) == 0 {
		return nil
	}
	return <-p.ErrChan
}

func (f *HandlerFunc) SetRets(ret ...interface{}) {
	f.Ret = ret
}

func (f *HandlerFunc) run() {
	inputs := make([]reflect.Value, len(f.Args))
	fT := reflect.ValueOf(f.F)
	for i := 0; i < len(f.Args); i++ {
		if f.Args[i] == nil {
			inputs[i] = reflect.Zero(fT.Type().In(i))
		} else {
			inputs[i] = reflect.ValueOf(f.Args[i])
		}
	}
	out := fT.Call(inputs)
	for i := 0; i < len(f.Ret); i++ {
		v := reflect.ValueOf(f.Ret[i])
		v.Elem().Set(out[i])
	}
}
