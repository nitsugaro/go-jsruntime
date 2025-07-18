package jsrun

import (
	"sync"

	"github.com/dop251/goja"
)

type ScriptManager struct {
	mu      sync.RWMutex
	scripts map[string]*goja.Program
	storage IScriptStorage
}

func NewScriptManager(storage IScriptStorage) *ScriptManager {
	sm := &ScriptManager{
		scripts: make(map[string]*goja.Program),
		storage: storage,
	}

	sm.storage.OnChange(func(name, code string) {
		sm.mustCompileScript(name, code)
	})

	return sm
}

func (sm *ScriptManager) mustCompileScript(name, code string) (*goja.Program, error) {
	prog, err := goja.Compile(name, code, false)
	if err != nil {
		return nil, err
	}

	sm.mu.Lock()
	sm.scripts[name] = prog
	sm.mu.Unlock()

	return prog, nil
}

func (sm *ScriptManager) DeleteFromCache(name string) {
	sm.mu.Lock()
	delete(sm.scripts, name)
	sm.mu.Unlock()
}

/* Compiles script and cached it */
func (sm *ScriptManager) CompileScript(name, code string) (*goja.Program, error) {
	sm.mu.RLock()
	if prog, ok := sm.scripts[name]; ok {
		sm.mu.RUnlock()
		return prog, nil
	}
	sm.mu.RUnlock()
	return sm.mustCompileScript(name, code)
}

func (sm *ScriptManager) ExecuteWithBindings(program *goja.Program, bindings map[string]interface{}) (goja.Value, error) {
	vm := goja.New()
	for k, v := range bindings {
		vm.Set(k, v)
	}

	vm.Set("require", sm.requireFunc(vm, bindings, make(map[string]goja.Value)))
	vm.Set("executeAsync", func(funcName string, bindings map[string]interface{}) {
		go sm.executeAsync(funcName, bindings)
	})
	scriptCallbacks := NewScriptCallbacks()
	vm.Set("callbacks", scriptCallbacks)

	val, err := vm.RunProgram(program)

	go scriptCallbacks.executeCallbacks()

	return val, err
}
