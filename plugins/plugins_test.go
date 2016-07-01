package plugins

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/modules"
)

func TestPlugins(t *testing.T) {
	viper.SetConfigFile("../config/test.yml")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	plugins := GetPlugins()
	if plugins == nil {
		t.Fail()
	}

	if !addCredentialsToRedis() {
		t.Fail()
	}
	plugins.SetupPlugins()

	payload := `{"payload": {"message": "register", "username": "username", "password": "userpass"}}`
	_, err := plugins.ExecutePlugin(payload, "mqttbot/acl/test", "register_user")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func addCredentialsToRedis() bool {
	user := viper.GetString("mqttserver.user")
	pass := viper.GetString("mqttserver.pass")
	hash := modules.GenHash(pass)
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
