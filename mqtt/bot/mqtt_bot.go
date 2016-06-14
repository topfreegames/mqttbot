package bot

import (
	"regexp"
	"strings"
	"sync"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/topfreegames/mqttbot/mqtt"
	"github.com/topfreegames/mqttbot/plugins"
)

type PluginMapping struct {
	Plugin         string
	MessagePattern string
}

type Subscription struct {
	Topic          string
	Qos            int
	PluginMappings []*PluginMapping
}

type MqttBot struct {
	Plugins       *plugins.Plugins
	Subscriptions []*Subscription
	Client        *mqttclient.MqttClient
}

var mqttBot *MqttBot
var once sync.Once

var h mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	for _, subscription := range mqttBot.Subscriptions {
		if RouteIncludesTopic(strings.Split(subscription.Topic, "/"), strings.Split(msg.Topic(), "/")) {
			for _, pluginMapping := range subscription.PluginMappings {
				match, _ := regexp.Match(pluginMapping.MessagePattern, msg.Payload())
				if match {
					mqttBot.Plugins.ExecutePlugin(string(msg.Payload()[:]), pluginMapping.Plugin)
					//trigger here
				}
			}
		}
	}
}

func GetMqttBot() *MqttBot {
	once.Do(func() {
		mqttBot = &MqttBot{}
		mqttBot.Client = mqttclient.GetMqttClient()
		mqttBot.setupPlugins()
	})
	return mqttBot
}

func (b *MqttBot) setupPlugins() {
	b.Plugins = plugins.GetPlugins()
	b.Plugins.SetupPlugins()
}

func (b *MqttBot) StartBot() {
	subscriptions := viper.Get("mqttserver.subscriptionRequests").([]interface{})
	client := b.Client.MqttClient
	b.Subscriptions = []*Subscription{}
	for _, s := range subscriptions {
		sMap := s.(map[interface{}]interface{})
		qos := sMap[string("qos")].(int)
		topic := sMap[string("topic")].(string)
		pluginMapping := sMap[string("plugins")].([]interface{})
		subscriptionNow := &Subscription{
			Topic:          topic,
			Qos:            qos,
			PluginMappings: []*PluginMapping{},
		}
		for _, p := range pluginMapping {
			pMap := p.(map[interface{}]interface{})
			subscriptionNow.PluginMappings = append(subscriptionNow.PluginMappings, &PluginMapping{
				Plugin:         pMap[string("plugin")].(string),
				MessagePattern: pMap[string("messagePattern")].(string),
			})
		}
		if token := client.Subscribe(topic, uint8(qos), h); token.Wait() && token.Error() != nil {
			logger.Logger.Fatal(token.Error())
		}
		b.Subscriptions = append(b.Subscriptions, subscriptionNow)
	}
	logger.Logger.Debug("Successfully subscribed to mqtt topics matching patterns /chat/# and /bot/history/#")
}
