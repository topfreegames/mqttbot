package plugins

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/modules"
)

func TestPlugins(t *testing.T) {
	config := viper.New()
	config.SetConfigFile("../config/test.yml")
	config.AutomaticEnv()
	config.ReadInConfig()
	plugins := GetPlugins(config)
	if plugins == nil {
		t.Fail()
	}

	if !addCredentialsToRedis(config) {
		t.Fail()
	}
	plugins.SetupPlugins()

	payload := `{"payload": {"message": "register", "username": "username", "password": "userpass"}}`
	_, err := plugins.ExecutePlugin(payload, "mqttbot/acl/test", "register_user")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func addCredentialsToRedis(config *viper.Viper) bool {
	user := config.GetString("mqttserver.user")
	pass := config.GetString("mqttserver.pass")
	hash := modules.GenHash(pass)
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
