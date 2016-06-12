package app

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"os"
)

type MqttClient struct {
	MqttServerHost string
	MqttServerPort int
	ConfigPath     string
	Config         *viper.Viper
	MqttClient     *client.Client
	ChatHandler    *ChatHandler
}

func GetMqttClient(mqttServerHost string, mqttServerPort int) *MqttClient {
	mqttClient := &MqttClient{
		MqttServerHost: mqttServerHost,
		MqttServerPort: mqttServerPort,
	}
	mqttClient.ChatHandler = GetChatHandler()
	return mqttClient
}

func (m *MqttClient) Start() {
	logger.Logger.Debug("Initializing mqtt client")
	m.MqttClient = client.New(&client.Options{
		ErrorHandler: func(err error) {
			logger.Logger.Error(err)
		},
	})

	c := m.MqttClient

	defer c.Terminate()
	err := c.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  fmt.Sprintf("%s:%d", m.MqttServerHost, m.MqttServerPort),
		ClientID: []byte("mqttbot-client"),
	})

	if err != nil {
		logger.Logger.Error("Error connecting to mqtt server! err:", err)
		os.Exit(1)
	}
	logger.Logger.Info(fmt.Sprintf("Successfully connected to mqtt server at %s:%d!", m.MqttServerHost, m.MqttServerPort))

	err = c.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("/chat/#"),
				QoS:         mqtt.QoS2,
				Handler: func(topicName, message []byte) {
				},
			},
			&client.SubReq{
				TopicFilter: []byte("/mqttbot/#"),
				QoS:         mqtt.QoS2,
				Handler: func(topicName, message []byte) {
					logger.Logger.Debug(fmt.Sprintf("Bot received message: %s, on topic: %s, serving history", string(message), string(topicName)))
				},
			},
		},
	})

	if err != nil {
		logger.Logger.Error("Error subscribing to mqtt topics! err:", err)
		os.Exit(1)
	}

	logger.Logger.Debug("Successfully subscribed to mqtt topics matching patterns /chat/# and /bot/history/#")

}
