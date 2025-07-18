package jsrun

import (
	"errors"

	"github.com/nitsugaro/go-nstore"
)

type IScriptStorage interface {
	GetSourceCode(name string, typ string) (string, error)
	OnChange(callback func(name, code string))
}

type ScriptStorage struct {
	*nstore.NStorage[*Script]

	callbacks []func(name, code string)
}

func (ss *ScriptStorage) OnChange(callback func(name, code string)) {
	if ss.callbacks == nil {
		ss.callbacks = []func(name, code string){}
	}

	ss.callbacks = append(ss.callbacks, callback)
}

func (ss *ScriptStorage) GetSourceCode(name string, typ string) (string, error) {
	if results, total := ss.Query(func(t *Script) bool { return t.Name == name && t.Type == typ }, 1); total == 1 {
		return results[0].GetRawCode()
	}

	return "", errors.New("not found")
}

func (ss *ScriptStorage) Save(script *Script) error {
	err := ss.NStorage.Save(script)
	if err != nil {
		return err
	}

	for _, cb := range ss.callbacks {
		if code, err := script.GetRawCode(); err != nil {
			cb(script.Name, code)
		}
	}

	return nil
}
