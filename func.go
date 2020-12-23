package pattern

import (
	"errors"
	"reflect"
)

type HandlerFunc struct {
	F    interface{}
	Args []interface{}
	Ret  []interface{}
}

func NewHandlerFunc(f interface{}, args ...interface{}) *HandlerFunc {
	handle := HandlerFunc{
		F:    f,
		Args: args,
		Ret:  []interface{}{},
	}
	return &handle
}
func (f *HandlerFunc) SetRets(ret ...interface{}) *HandlerFunc {
	f.Ret = ret
	return f
}

func (f *HandlerFunc) Run() {
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

func (f *HandlerFunc) Check() error {
	for _, ret := range f.Ret {
		if reflect.ValueOf(ret).Type().Kind() != reflect.Ptr {
			return errors.New("the ret is not ptr")
		}
	}
	fT := reflect.ValueOf(f.F).Type()
	if fT.Kind() != reflect.Func {
		return errors.New("the func is not function")
	}
	if fT.NumOut() != len(f.Ret) {
		return errors.New("the func return is not valid")
	}
	return nil
}
