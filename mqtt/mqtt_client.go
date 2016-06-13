package mqttclient

import (
	"fmt"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
)

type MqttClient struct {
	MqttServerHost string
	MqttServerPort int
	ConfigPath     string
	MqttClient     mqtt.Client
}

var Client *MqttClient
var once sync.Once

func GetMqttClient() *MqttClient {
	once.Do(func() {
		Client = &MqttClient{}
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

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", mc.MqttServerHost, mc.MqttServerPort)).SetClientID("mqttbot")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	mc.MqttClient = mqtt.NewClient(opts)

	c := mc.MqttClient

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Logger.Fatal(token.Error())
	}

	logger.Logger.Info(fmt.Sprintf("Successfully connected to mqtt server at %s:%d!", mc.MqttServerHost, mc.MqttServerPort))
}
