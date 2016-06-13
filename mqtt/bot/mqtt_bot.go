package bot

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/topfreegames/mqttbot/mqtt"
	"github.com/topfreegames/mqttbot/plugins"
)

type MqttBot struct {
	Plugins         *plugins.Plugins
	PluginsMappings []map[string]string
	Client          *mqttclient.MqttClient
}

var h mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
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

func (b *MqttBot) StartBot() {
	subscriptions := viper.Get("mqttserver.subscriptionRequests").([]interface{})
	client := b.Client.MqttClient
	for _, s := range subscriptions {
		sMap := s.(map[interface{}]interface{})
		qos := sMap[string("qos")].(int)
		if token := client.Subscribe(sMap[string("topic")].(string), uint8(qos), h); token.Wait() && token.Error() != nil {
			logger.Logger.Fatal(token.Error())
		}
	}

	logger.Logger.Debug("Successfully subscribed to mqtt topics matching patterns /chat/# and /bot/history/#")
}
