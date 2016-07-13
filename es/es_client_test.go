package es

import (
	"testing"

	"github.com/spf13/viper"
)

func TestES(t *testing.T) {
	viper.SetConfigFile("../config/test.yaml")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	client := GetESClient()
	if client == nil {
		t.Fail()
	}
}
