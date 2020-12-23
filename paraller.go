// this paraller model is used to execute functions in parallel
package pattern

import (
	"fmt"
	"sync"
)

type Paraller struct {
	wg       sync.WaitGroup
	handlers []*HandlerFunc
	ErrChan  chan error
}

func NewParaller() *Paraller {
	return &Paraller{}
}

func (p *Paraller) Add(handle *HandlerFunc) {
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
			go func(f HandlerFunc, ch chan error) {
				defer p.wg.Done()
				defer func() {
					if err := recover(); err != nil {
						ch <- fmt.Errorf("[recover]error: %+v", err)
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
