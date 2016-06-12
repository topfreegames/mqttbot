package app

import (
	"github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"index": Index,
	"query": Query,
}

func IndexMessage(L *lua.LState) int {
	return 0
}

func QueryMessages(L *lua.LState) int {
	return 0
}
