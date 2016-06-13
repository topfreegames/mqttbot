package bot

import (
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/topfreegames/mqttbot/mqtt"
	"github.com/topfreegames/mqttbot/plugins"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type MqttBot struct {
	Plugins         *plugins.Plugins
	PluginsMappings []map[string]string
	Client          *mqttclient.MqttClient
}

func GetMqttBot() *MqttBot {
	mqttBot := &MqttBot{}
	mqttBot.Client = mqttclient.GetMqttClient()
	mqttBot.setupPlugins()
	return mqttBot
}

func (b *MqttBot) setupPlugins() {
	b.Plugins = plugins.GetPlugins()
	b.Plugins.SetupPlugins()
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

func (b *MqttBot) StartBot() {
	subscriptions := viper.Get("mqttserver.subscriptionRequests").([]interface{})
	subscriptionsOptions := &client.SubscribeOptions{SubReqs: []*client.SubReq{}}
	for _, s := range subscriptions {
		sMap := s.(map[interface{}]interface{})
		//pluginsMap := sMap[string("plugins")].([]interface{})

		subscriptionReq := &client.SubReq{
			TopicFilter: []byte(sMap[string("topic")].(string)),
			QoS:         getQosFromInt(sMap[string("qos")].(int)),
			Handler: func(topicName, message []byte) {
			},
		}
		subscriptionsOptions.SubReqs = append(subscriptionsOptions.SubReqs, subscriptionReq)
	}

	logger.Logger.Debug("Successfully subscribed to mqtt topics matching patterns /chat/# and /bot/history/#")
}
