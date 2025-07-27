package script

import (
	"path/filepath"

	"github.com/yuin/gopher-lua"

	"github.com/kkkunny/pokemon/src/config"
)

func LoadScriptFile(name string) (*lua.LState, error) {
	l := lua.NewState()
	lf, err := l.LoadFile(filepath.Join(config.ScriptsPath, name) + ".lua")
	if err != nil {
		return nil, err
	}
	l.Push(lf)
	return l, nil
}
