package mqttclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
)

// MQTTClient contains the data needed to connect the client
type MQTTClient struct {
	MQTTServerHost   string
	MQTTServerPort   int
	MQTTClient       mqtt.Client
	Heartbeat        *Heartbeat
	OnConnectHandler mqtt.OnConnectHandler
}

var client *MQTTClient
var once sync.Once

//ResetMQTTClient resets once
func ResetMQTTClient() {
	once = sync.Once{}
	client = nil
}

// GetMQTTClient creates the mqttclient and returns it
func GetMQTTClient(onConnectHandler mqtt.OnConnectHandler) *MQTTClient {
	once.Do(func() {
		client = &MQTTClient{}
		client.configure(onConnectHandler)
	})
	return client
}

func (mc *MQTTClient) hasConnected(client mqtt.Client) {
	if mc.OnConnectHandler != nil {
		mc.OnConnectHandler(client)
	}

	// mc.Heartbeat = &Heartbeat{
	// 	Topic:             uuid.NewV4().String(),
	// 	Client:            mc,
	// 	OnHeartbeatMissed: mc.onHeartbeatMissed,
	// }
	// mc.Heartbeat.Start()
}

func (mc *MQTTClient) configure(onConnectHandler mqtt.OnConnectHandler) {
	mc.OnConnectHandler = onConnectHandler
	mc.setConfigurationDefaults()
	mc.configureClient()
	mc.start()
}

// func (mc *MQTTClient) onHeartbeatMissed(err error) {
// 	logger.Logger.Info("Heartbeat missed")
// 	if mc.MQTTClient.IsConnected() {
// 		logger.Logger.Info("Heartbeat missed: disconnecting")
// 		mc.MQTTClient.Disconnect(0)
// 	}
// 	logger.Logger.Info("Heartbeat missed: connecting again")
// 	mc.start()
// }

func (mc *MQTTClient) setConfigurationDefaults() {
	viper.SetDefault("mqttserver.host", "localhost")
	viper.SetDefault("mqttserver.port", 1883)
	viper.SetDefault("mqttserver.user", "admin")
	viper.SetDefault("mqttserver.pass", "admin")
	viper.SetDefault("mqttserver.subscriptions", []map[string]string{})
	viper.SetDefault("mqttserver.ca_cert_file", "")
}

func (mc *MQTTClient) configureClient() {
	mc.MQTTServerHost = viper.GetString("mqttserver.host")
	mc.MQTTServerPort = viper.GetInt("mqttserver.port")
}

func (mc *MQTTClient) start() {
	logger.Logger.Debug("Initializing mqtt client")

	useTLS := viper.GetBool("mqttserver.usetls")
	protocol := "tcp"
	if useTLS {
		protocol = "ssl"
	}

	opts := mqtt.NewClientOptions().AddBroker(
		fmt.Sprintf("%s://%s:%d", protocol, mc.MQTTServerHost, mc.MQTTServerPort),
	).SetClientID("mqttbot")

	if useTLS {
		logger.Logger.Info("mqttclient using tls")
		certpool := x509.NewCertPool()
		if viper.GetString("mqttserver.ca_cert_file") != "" {
			pemCerts, err := ioutil.ReadFile(viper.GetString("mqttserver.ca_cert_file"))
			if err == nil {
				certpool.AppendCertsFromPEM(pemCerts)
			} else {
				logger.Logger.Error(err.Error())
			}
		}
		tlsConfig := &tls.Config{
			InsecureSkipVerify: viper.GetBool("mqttserver.insecure_tls"),
			ClientAuth:         tls.NoClientCert,
			RootCAs:            certpool,
		}
		opts.SetTLSConfig(tlsConfig)
	}

	opts.SetUsername(viper.GetString("mqttserver.user"))
	opts.SetPassword(viper.GetString("mqttserver.pass"))
	opts.SetKeepAlive(15 * time.Second)
	opts.SetPingTimeout(60 * time.Second)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetOnConnectHandler(mc.hasConnected)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetAutoReconnect(true)
	mc.MQTTClient = mqtt.NewClient(opts)

	c := mc.MQTTClient

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Logger.Fatal(token.Error())
	}

	logger.Logger.Info(fmt.Sprintf(
		"Successfully connected to mqtt server at %s:%d!",
		mc.MQTTServerHost, mc.MQTTServerPort,
	))
}
