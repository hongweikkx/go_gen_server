// actor model.
package actor

import (
	"errors"
	"fmt"
	"time"

	pattern "github.com/hongweikkx/go_pattern"
)

type Mod struct {
	mailBox *MailBox
}

type MailBox struct {
	callCh     chan *pattern.HandlerFunc
	castCh     chan *pattern.HandlerFunc
	exitCh     chan string
	callDoneCh chan bool
}

// castNum : 异步channel的大小
func Start(castNum int) *Mod {
	mod := &Mod{
		mailBox: &MailBox{
			callCh:     make(chan *pattern.HandlerFunc, 1),
			castCh:     make(chan *pattern.HandlerFunc, castNum),
			exitCh:     make(chan string),
			callDoneCh: make(chan bool),
		}}
	go mod.mailBox.doSpawn()
	return mod
}

func (mod *Mod) Call(msg *pattern.HandlerFunc) error {
	fmt.Println(time.Now().Unix())
	mod.mailBox.callCh <- msg
	fmt.Println(time.Now().Unix(), mod.mailBox.callDoneCh)
	select {
	case <-mod.mailBox.callDoneCh:
		return nil
	case <-time.After(5 * time.Second):
		fmt.Println(time.Now().Unix())
		return errors.New("call timeout")
	}
}

func (mod *Mod) Cast(msg *pattern.HandlerFunc) {
	mod.mailBox.castCh <- msg
}

func (mod *Mod) Stop(msg string) {
	mod.mailBox.exitCh <- msg
}

func (mailBox *MailBox) doSpawn() {
	mailBox.loop()
	mailBox.terminate()
	return
}

func (mailBox *MailBox) loop() string {
	for {
		select {
		case callMsg := <-mailBox.callCh:
			callMsg.Run()
			mailBox.callDoneCh <- true
		case castMsg := <-mailBox.castCh:
			castMsg.Run()
		case exitReason := <-mailBox.exitCh:
			return exitReason
		}
	}
}

func (mailBox *MailBox) terminate() {
	close(mailBox.callCh)
	close(mailBox.callDoneCh)
	close(mailBox.castCh)
	close(mailBox.exitCh)
}
