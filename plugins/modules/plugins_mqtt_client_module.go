package modules

import (
	"github.com/kataras/iris/logger"
	"github.com/topfreegames/mqttbot/mqtt"
	"github.com/yuin/gopher-lua"
)

var mqttClient *mqtt.MqttClient

func MqttClientModuleLoader(L *lua.LState) int {
	configureMqttModule()
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"sendMessage": SendMessage,
}

func configureMqttModule() {
	mqttClient = mqtt.GetMqttClient()
}

func SendMessage(L *lua.LState) int {
	logger.Logger.Debug("mqttClientModule SendMessage called")
	return 0
}
