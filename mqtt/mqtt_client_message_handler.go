package mqtt

import (
	"github.com/topfreegames/mqttbot/plugins"
)

type ChatHandler struct {
	Plugins *plugins.Plugins
}

type ChatMessage struct {
	TopicName []byte
	Payload   []byte
}

func GetChatHandler() *ChatHandler {
	handler := &ChatHandler{
		Plugins: plugins.GetPlugins(),
	}
	p := handler.Plugins
	p.SetupPlugins()
	return handler
}
