// this paraller model is used to execute functions in parallel
package paraller

import (
	"errors"
	"fmt"
	"sync"

	pattern "github.com/hongweikkx/go_pattern"
)

type Paraller struct {
	wg       sync.WaitGroup
	handlers []*pattern.HandlerFunc
	ErrChan  chan error
}

func NewParaller() *Paraller {
	return &Paraller{}
}

func (p *Paraller) Add(handle *pattern.HandlerFunc) {
	p.handlers = append(p.handlers, handle)
}

func (p *Paraller) Run() error {
	if err := p.checkPre(); err != nil {
		return err
	}
	if hl := len(p.handlers); hl > 0 {
		p.ErrChan = make(chan error, hl)
		p.wg.Add(len(p.handlers))
		for k := range p.handlers {
			go func(f pattern.HandlerFunc, ch chan error) {
				defer p.wg.Done()
				defer func() {
					if err := recover(); err != nil {
						ch <- errors.New(fmt.Sprintf("[recover]error: %+v", err))
					}
				}()
				f.Run()
			}(*p.handlers[k], p.ErrChan)
		}
		p.wg.Wait()
	}
	return p.errorRet()
}

func (p *Paraller) checkPre() error {
	for _, handler := range p.handlers {
		err := handler.Check()
		if err != nil {
			return err
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
