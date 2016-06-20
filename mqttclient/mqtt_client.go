package mqttclient

import (
	"fmt"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
)

// MqttClient contains the data needed to connect the client
type MqttClient struct {
	MqttServerHost string
	MqttServerPort int
	Config         *viper.Viper
	MqttClient     mqtt.Client
}

var client *MqttClient
var once sync.Once

// GetMqttClient creates the mqttclient and returns it
func GetMqttClient(config *viper.Viper, onConnectHandler mqtt.OnConnectHandler) *MqttClient {
	once.Do(func() {
		client = &MqttClient{Config: config}
		client.configure(onConnectHandler)
	})
	return client
}

func (mc *MqttClient) configure(onConnectHandler mqtt.OnConnectHandler) {
	mc.setConfigurationDefaults()
	mc.configureClient()
	mc.start(onConnectHandler)
}

func (mc *MqttClient) setConfigurationDefaults() {
	mc.Config.SetDefault("mqttserver.host", "localhost")
	mc.Config.SetDefault("mqttserver.port", 1883)
	mc.Config.SetDefault("mqttserver.user", "admin")
	mc.Config.SetDefault("mqttserver.pass", "admin")
	mc.Config.SetDefault("mqttserver.subscriptions", []map[string]string{})
}

func (mc *MqttClient) configureClient() {
	mc.MqttServerHost = mc.Config.GetString("mqttserver.host")
	mc.MqttServerPort = mc.Config.GetInt("mqttserver.port")
}

func (mc *MqttClient) start(onConnectHandler mqtt.OnConnectHandler) {
	logger.Logger.Debug("Initializing mqtt client")

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", mc.MqttServerHost, mc.MqttServerPort)).SetClientID("mqttbot")
	opts.SetUsername(mc.Config.GetString("mqttserver.user"))
	opts.SetPassword(mc.Config.GetString("mqttserver.pass"))
	opts.SetKeepAlive(3 * time.Second)
	opts.SetPingTimeout(5 * time.Second)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetOnConnectHandler(onConnectHandler)
	mc.MqttClient = mqtt.NewClient(opts)

	c := mc.MqttClient

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Logger.Fatal(token.Error())
	}

	logger.Logger.Info(fmt.Sprintf("Successfully connected to mqtt server at %s:%d!",
		mc.MqttServerHost, mc.MqttServerPort))
}
