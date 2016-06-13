package mqtt

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type MqttClient struct {
	MqttServerHost string
	MqttServerPort int
	ConfigPath     string
	MqttClient     *client.Client
	Config         *viper.Viper
	ChatHandler    *ChatHandler
}

var Client *MqttClient
var once sync.Once

func GetMqttClient() *MqttClient {
	once.Do(func() {
		Client := &MqttClient{
			Config: viper.New(),
		}
		Client.ChatHandler = GetChatHandler()
		Client.configure()
	})
	return Client
}

func (c *MqttClient) configure() {
	c.setConfigurationDefaults()
	c.loadConfiguration()
	c.configureClient()
	c.start()
}

func (c *MqttClient) setConfigurationDefaults() {
	c.Config.SetDefault("mqttserver.host", "localhost")
	c.Config.SetDefault("mqttserver.port", 1883)
}

func (c *MqttClient) loadConfiguration() {
	c.Config.SetConfigFile("./config/mqtt.yaml")
	c.Config.SetEnvPrefix("mqttbot.mqtt")
	c.Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.Config.AutomaticEnv()

	if err := c.Config.ReadInConfig(); err == nil {
		logger.Logger.Debug(fmt.Sprintf("Using config file: %s", c.Config.ConfigFileUsed()))
	}
}

func (c *MqttClient) configureClient() {
	c.MqttServerHost = c.Config.GetString("mqttserver.host")
	c.MqttServerPort = c.Config.GetInt("mqttserver.port")
}

func (m *MqttClient) start() {
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
