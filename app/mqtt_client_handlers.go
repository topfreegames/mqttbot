package app

import (
	"fmt"
	"github.com/topfreegames/mqttbot/models"
	"gopkg.in/olivere/elastic.v3"
)

type ChatHandler struct {
	ESClient *elastic.Client
}

type ChatMessage struct {
	TopicName []byte
	Payload   []byte
}

func GetChatHandler(esClient *elastic.Client) *ChatHandler {
	handler := &ChatHandler{
		ESClient: esClient,
	}
	return handler
}

func (c *ChatHandler) PersistMessage(chatMessage *ChatMessage) {
	Logger.Debug(fmt.Sprintf("Persisting message: %s on topic: %s", chatMessage.Payload, chatMessage.TopicName))
	message := &models.Message{
		Message: fmt.Sprintf("%s", chatMessage.Payload),
		Topic:   fmt.Sprintf("%s", chatMessage.TopicName),
	}

	err := message.Index(c.ESClient)
	if err != nil {
		Logger.Error("Error persisting message into chat index! err:", err)
	} else {
		Logger.Debug("Message successfully persisted into chat index!")
	}
}

//func (c *ChatHandler) HandleBotMessage(chatMessage *ChatMessage) {

//}
