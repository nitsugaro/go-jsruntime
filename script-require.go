package jsrun

import (
	"fmt"

	"github.com/dop251/goja"
)

func (sm *ScriptManager) requireFunc(vm *goja.Runtime, bindings map[string]interface{}, loadedModules map[string]goja.Value) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		moduleName := call.Argument(0).String()
		if val, ok := loadedModules[moduleName]; ok {
			return val
		}

		sm.mu.RLock()
		cached, ok := sm.scripts[moduleName]
		sm.mu.RUnlock()
		if !ok {
			panic(vm.ToValue(vm.NewGoError(fmt.Errorf("module %q not found", moduleName))))
		}

		for k, v := range bindings {
			vm.Set(k, v)
		}

		vm.Set("require", sm.requireFunc(vm, bindings, loadedModules))
		vm.Set("executeAsync", func(funcName string, bindings map[string]interface{}) {
			go sm.executeAsync(funcName, bindings)
		})

		exports := vm.NewObject()
		module := vm.NewObject()
		module.Set("exports", exports)

		vm.Set("module", module)
		vm.Set("exports", exports)

		fnVal, err := vm.RunProgram(cached)
		if err != nil {
			panic(vm.ToValue(vm.NewGoError(fmt.Errorf("error executing module %q: %w", moduleName, err))))
		}

		fn, ok := goja.AssertFunction(fnVal)
		if !ok {
			panic(vm.ToValue(vm.NewGoError(fmt.Errorf("module doesn't %q export any function", moduleName))))
		}

		_, err = fn(goja.Undefined(), exports, module)
		if err != nil {
			panic(vm.ToValue(vm.NewGoError((err))))
		}

		result := module.Get("exports")

		loadedModules[moduleName] = result

		return result
	}
}
