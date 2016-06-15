package modules

import (
	"fmt"

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
	"send_message": SendMessage,
}

func configureMqttModule() {
	mqttClient = mqttclient.GetMqttClient().MqttClient
}

func SendMessage(L *lua.LState) int {
	topic := L.Get(-4)
	qos := L.Get(-3)
	retained := L.Get(-2)
	payload := L.Get(-1)
	L.Pop(4)
	logger.Logger.Debug(fmt.Sprintf("mqttclient_module send message topic: %s, payload: %s, qos: %d, retained: %s", topic, payload, qos, retained))
	if token := mqttClient.Publish(topic.String(), byte(qos.(lua.LNumber)), bool(retained.(lua.LBool)), payload.String()); token.Wait() && token.Error() != nil {
		logger.Logger.Error(token.Error())
		L.Push(lua.LString(fmt.Sprintf("%s", token.Error())))
		L.Push(L.ToNumber(1))
		return 2
	}
	L.Push(lua.LNil)
	L.Push(L.ToNumber(0))
	return 2
}
