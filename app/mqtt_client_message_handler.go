package app

import ()

type ChatHandler struct {
	Plugins *Plugins
}

type ChatMessage struct {
	TopicName []byte
	Payload   []byte
}

func GetChatHandler() *ChatHandler {
	handler := &ChatHandler{
		Plugins: GetPlugins(),
	}
	p := handler.Plugins
	p.SetupPlugins()
	return handler
}
