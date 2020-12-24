// this parallel model is used to execute functions in parallel
package pattern

import (
	"fmt"
	"sync"
)

type Parallel struct {
	wg       sync.WaitGroup
	handlers []*HandlerFunc
	ErrChan  chan error
}

func NewParallel() *Parallel {
	return &Parallel{}
}

func (p *Parallel) Add(handle *HandlerFunc) {
	p.handlers = append(p.handlers, handle)
}

func (p *Parallel) Run() error {
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

func (p *Parallel) checkPre() error {
	for _, handler := range p.handlers {
		err := handler.Check()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parallel) errorRet() error {
	if len(p.ErrChan) == 0 {
		return nil
	}
	return <-p.ErrChan
}
