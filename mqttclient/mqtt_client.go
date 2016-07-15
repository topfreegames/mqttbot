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
	viper.SetDefault("mqttserver.ca_cert_file", "")
}

func (mc *MqttClient) configureClient() {
	mc.MqttServerHost = viper.GetString("mqttserver.host")
	mc.MqttServerPort = viper.GetInt("mqttserver.port")
}

func (mc *MqttClient) start(onConnectHandler mqtt.OnConnectHandler) {
	logger.Logger.Debug("Initializing mqtt client")

	useTls := viper.GetBool("mqttserver.usetls")
	protocol := "tcp"
	if useTls {
		protocol = "ssl"
	}

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("%s://%s:%d", protocol, mc.MqttServerHost, mc.MqttServerPort)).SetClientID("mqttbot")

	if useTls {
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
		tlsConfig := &tls.Config{InsecureSkipVerify: viper.GetBool("mqttserver.insecure_tls"), ClientAuth: tls.NoClientCert, RootCAs: certpool}
		opts.SetTLSConfig(tlsConfig)
	}

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
