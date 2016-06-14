package modules

import (
	"github.com/topfreegames/mqttbot/es"
	"github.com/yuin/gopher-lua"
	"gopkg.in/olivere/elastic.v3"
)

var ESClient *elastic.Client

func PersistenceModuleLoader(L *lua.LState) int {
	configurePersistenceModule()
	mod := L.SetFuncs(L.NewTable(), persistenceModuleExports)
	L.Push(mod)
	return 1
}

var persistenceModuleExports = map[string]lua.LGFunction{
	"index_message":  IndexMessage,
	"query_messages": QueryMessages,
}

func configurePersistenceModule() {
	ESClient = es.GetESClient()
}

func IndexMessage(L *lua.LState) int {
	return 0
}

func QueryMessages(L *lua.LState) int {
	return 0
}
