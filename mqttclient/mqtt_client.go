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
	MqttClient     mqtt.Client
}

var client *MqttClient
var once sync.Once

// GetMqttClient creates the mqttclient and returns it
func GetMqttClient(onConnectHandler mqtt.OnConnectHandler) *MqttClient {
	once.Do(func() {
		client = &MqttClient{}
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
	viper.SetDefault("mqttserver.host", "localhost")
	viper.SetDefault("mqttserver.port", 1883)
	viper.SetDefault("mqttserver.user", "admin")
	viper.SetDefault("mqttserver.pass", "admin")
	viper.SetDefault("mqttserver.subscriptions", []map[string]string{})
}

func (mc *MqttClient) configureClient() {
	mc.MqttServerHost = viper.GetString("mqttserver.host")
	mc.MqttServerPort = viper.GetInt("mqttserver.port")
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

	logger.Logger.Info(fmt.Sprintf("Successfully connected to mqtt server at %s:%d!",
		mc.MqttServerHost, mc.MqttServerPort))
}
