package modules

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/topfreegames/mqttbot/mqtt"
	"github.com/yuin/gopher-lua"
)

var mqttClient mqtt.Client

func MqttClientModuleLoader(L *lua.LState) int {
	configureMqttModule()
	mod := L.SetFuncs(L.NewTable(), mqttClientModuleExports)
	L.Push(mod)
	return 1
}

var mqttClientModuleExports = map[string]lua.LGFunction{
	"sendMessage": SendMessage,
}

func configureMqttModule() {
	mqttClient = mqttclient.GetMqttClient().MqttClient
}

func SendMessage(L *lua.LState) int {
	logger.Logger.Debug("mqttClientModule SendMessage called")
	return 0
}
