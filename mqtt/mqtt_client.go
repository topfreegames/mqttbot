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
	MqttClient     mqtt.Client
}

var Client *MqttClient
var once sync.Once

func GetMqttClient(onConnectHandler mqtt.OnConnectHandler) *MqttClient {
	once.Do(func() {
		Client = &MqttClient{}
		Client.configure(onConnectHandler)
	})
	return Client
}

func (c *MqttClient) configure(onConnectHandler mqtt.OnConnectHandler) {
	c.setConfigurationDefaults()
	c.configureClient()
	c.start(onConnectHandler)
}

func (c *MqttClient) setConfigurationDefaults() {
	viper.SetDefault("mqttserver.host", "localhost")
	viper.SetDefault("mqttserver.port", 1883)
	viper.SetDefault("mqttserver.user", "admin")
	viper.SetDefault("mqttserver.pass", "admin")
	viper.SetDefault("mqttserver.subscriptions", []map[string]string{})
}

func (c *MqttClient) configureClient() {
	c.MqttServerHost = viper.GetString("mqttserver.host")
	c.MqttServerPort = viper.GetInt("mqttserver.port")
}

func (mc *MqttClient) start(onConnectHandler mqtt.OnConnectHandler) {
	logger.Logger.Debug("Initializing mqtt client")

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", mc.MqttServerHost, mc.MqttServerPort)).SetClientID("mqttbot")
	opts.SetUsername(viper.GetString("mqttserver.user"))
	opts.SetPassword(viper.GetString("mqttserver.pass"))
	opts.SetKeepAlive(3 * time.Second)
	opts.SetPingTimeout(5 * time.Second)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetOnConnectHandler(onConnectHandler)
	mc.MqttClient = mqtt.NewClient(opts)

	c := mc.MqttClient

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Logger.Fatal(token.Error())
	}

	logger.Logger.Info(fmt.Sprintf("Successfully connected to mqtt server at %s:%d!", mc.MqttServerHost, mc.MqttServerPort))
}
