package mqttclient

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/yosssi/gmq/mqtt/client"
)

type MqttClient struct {
	MqttServerHost string
	MqttServerPort int
	ConfigPath     string
	MqttClient     *client.Client
}

var Client *MqttClient
var once sync.Once

func GetMqttClient() *MqttClient {
	once.Do(func() {
		Client := &MqttClient{}
		Client.configure()
	})
	return Client
}

func (c *MqttClient) configure() {
	c.setConfigurationDefaults()
	c.configureClient()
	c.start()
}

func (c *MqttClient) setConfigurationDefaults() {
	viper.SetDefault("mqttserver.host", "localhost")
	viper.SetDefault("mqttserver.port", 1883)
	viper.SetDefault("mqttserver.subscriptions", []map[string]string{})
}

func (c *MqttClient) configureClient() {
	c.MqttServerHost = viper.GetString("mqttserver.host")
	c.MqttServerPort = viper.GetInt("mqttserver.port")
}

func (mc *MqttClient) start() {
	logger.Logger.Debug("Initializing mqtt client")
	mc.MqttClient = client.New(&client.Options{
		ErrorHandler: func(err error) {
			logger.Logger.Error(err)
		},
	})

	c := mc.MqttClient

	defer c.Terminate()
	err := c.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  fmt.Sprintf("%s:%d", mc.MqttServerHost, mc.MqttServerPort),
		ClientID: []byte("mqttbot-client"),
	})

	if err != nil {
		logger.Logger.Fatal("Error connecting to mqtt server! err:", err)
	}
	logger.Logger.Info(fmt.Sprintf("Successfully connected to mqtt server at %s:%d!", mc.MqttServerHost, mc.MqttServerPort))
}
