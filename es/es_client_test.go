package es

import (
	"testing"

	"github.com/spf13/viper"
)

func TestES(t *testing.T) {
	config := viper.New()
	config.SetConfigFile("../config/test.yml")
	config.AutomaticEnv()
	config.ReadInConfig()
	client := GetESClient(config)
	if client == nil {
		t.Fail()
	}
}
