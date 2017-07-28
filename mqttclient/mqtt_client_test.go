package mqttclient

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	"golang.org/x/crypto/pbkdf2"
)

func TestClient(t *testing.T) {
	ResetMQTTClient()
	viper.SetConfigFile("../config/test.yaml")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	if !addCredentialsToRedis() {
		t.Fail()
	}

	mqttClient := GetMQTTClient(nil)
	if mqttClient == nil {
		t.Fail()
	}
}

func genHash(pass string) string {
	bpass := []byte(pass)
	iterations := 901
	salt := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return ""
	}
	esalt := base64.StdEncoding.EncodeToString(salt)
	bhash := pbkdf2.Key(bpass, []byte(esalt), iterations, 24, sha256.New)
	ehash := base64.StdEncoding.EncodeToString(bhash)
	hash := fmt.Sprintf("PBKDF2$sha256$%d$%s$%s", iterations, esalt, ehash)
	return hash
}

func addCredentialsToRedis() bool {
	user := viper.GetString("mqttserver.user")
	pass := viper.GetString("mqttserver.pass")
	hash := genHash(pass)
	redisHost := viper.GetString("redis.host")
	redisPort := viper.GetInt("redis.port")
	redisPass := viper.GetString("redis.password")
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", redisHost, redisPort),
		redis.DialPassword(redisPass))
	if err != nil {
		return false
	}
	defer conn.Close()
	if _, err = conn.Do("SET", user, hash); err != nil {
		return false
	}
	return true
}

// func TestHeartbeat(t *testing.T) {
// 	mqttClient := GetMQTTClient(nil)
// 	if mqttClient == nil {
// 		t.Fail()
// 	}
//
// 	time.Sleep(100 * time.Millisecond)
//
// 	//Heartbeat not working
// 	if mqttClient.Heartbeat.LastHeartbeat.Unix() < time.Now().Add(-10*time.Second).Unix() {
// 		t.Fail()
// 	}
// }
//
// func TestHeartbeatReconnects(t *testing.T) {
// 	mqttClient := GetMQTTClient(nil)
// 	mqttClient.Heartbeat.MaxDurationMs = 400
// 	if mqttClient == nil {
// 		t.Fail()
// 	}
//
// 	time.Sleep(100 * time.Millisecond)
//
// 	//Heartbeat not working
// 	if mqttClient.Heartbeat.LastHeartbeat.Unix() < time.Now().Add(-10*time.Second).Unix() {
// 		t.Fail()
// 	}
//
// 	mqttClient.MQTTClient.Disconnect(0)
// 	if mqttClient.MQTTClient.IsConnected() {
// 		t.Fail()
// 	}
//
// 	time.Sleep(500 * time.Millisecond)
//
// 	if !mqttClient.MQTTClient.IsConnected() {
// 		t.Fatal("Should be connected!")
// 	}
// }
