package jsrun

import "github.com/dop251/goja"

type scriptCallbacks struct {
	callbacks []goja.Callable
}

func NewScriptCallbacks() *scriptCallbacks {
	return &scriptCallbacks{
		callbacks: []goja.Callable{},
	}
}

func (sc *scriptCallbacks) Add(cb goja.Callable) {
	sc.callbacks = append(sc.callbacks, cb)
}

func (sc *scriptCallbacks) executeCallbacks() {
	for _, cb := range sc.callbacks {
		_, err := cb(goja.Undefined())
		if err != nil {
		}
	}
}
