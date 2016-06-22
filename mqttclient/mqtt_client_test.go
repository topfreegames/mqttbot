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
	config := viper.New()
	config.SetConfigFile("../config/test.yml")
	config.AutomaticEnv()
	config.ReadInConfig()

	if !addCredentialsToRedis(config) {
		t.Fail()
	}

	mqttClient := GetMqttClient(config, nil)
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

func addCredentialsToRedis(config *viper.Viper) bool {
	user := config.GetString("mqttserver.user")
	pass := config.GetString("mqttserver.pass")
	hash := genHash(pass)
	redisHost := config.GetString("redis.host")
	redisPort := config.GetInt("redis.port")
	redisPass := config.GetString("redis.password")
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
