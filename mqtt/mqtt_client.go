package mqtt

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type MqttClient struct {
	MqttServerHost  string
	MqttServerPort  int
	ConfigPath      string
	MqttClient      *client.Client
	PluginsMappings []map[string]string
}

type subscriptionStruct struct {
	topic   string
	qos     int
	plugins []map[string]string
}

type subscriptionsStruct struct {
	Subscriptions []subscriptionStruct
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

func getQosFromInt(qosInt int) byte {
	switch qosInt {
	case 0:
		return mqtt.QoS0
	case 1:
		return mqtt.QoS1
	case 2:
		return mqtt.QoS2
	default:
		return mqtt.QoS2
	}
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

	subscriptions := viper.Get("mqttserver.subscriptions").([]interface{})
	subscriptionsOptions := &client.SubscribeOptions{SubReqs: []*client.SubReq{}}
	for _, s := range subscriptions {
		sMap := s.(map[interface{}]interface{})
		pluginsMap := sMap[string("plugins")].([]interface{})

		subscriptionReq := &client.SubReq{
			TopicFilter: []byte(sMap[string("topic")].(string)),
			QoS:         getQosFromInt(sMap[string("qos")].(int)),
			Handler: func(topicName, message []byte) {
			},
		}
		subscriptionsOptions.SubReqs = append(subscriptionsOptions.SubReqs, subscriptionReq)
	}

	if err != nil {
		logger.Logger.Error("Error subscribing to mqtt topics! err:", err)
		os.Exit(1)
	}

	logger.Logger.Debug("Successfully subscribed to mqtt topics matching patterns /chat/# and /bot/history/#")

}
