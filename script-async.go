package jsrun

import (
	"fmt"

	"github.com/dop251/goja"
)

func (sm *ScriptManager) executeAsync(funcName string, bindings map[string]interface{}) error {
	sm.mu.RLock()
	cached, ok := sm.scripts[funcName]
	sm.mu.RUnlock()

	if !ok {
		code, err := sm.storage.GetSourceCode(funcName, "async")
		if err != nil {
			return fmt.Errorf("module %q not found", funcName)
		}
		prog, err := sm.mustCompileScript(funcName, code)
		if err != nil {
			return fmt.Errorf("error compiling module %q: %w", funcName, err)
		}

		sm.mu.Lock()
		sm.scripts[funcName] = prog
		sm.mu.Unlock()

		cached = prog
	}

	vm := goja.New()
	for key, val := range bindings {
		vm.Set(key, val)
	}

	_, err := vm.RunProgram(cached)
	if err != nil {
		return err
	}

	return nil
}
