package jsrun

import (
	"github.com/nitsugaro/go-nstore"
)

func NewDefaultStorage(folder string) (*ScriptManager, *ScriptStorage) {
	var scriptStorage *ScriptStorage

	storage, err := nstore.New[*Script](folder)
	if err != nil {
		panic(err)
	}

	scriptStorage = &ScriptStorage{NStorage: storage}

	return NewScriptManager(scriptStorage), scriptStorage
}
