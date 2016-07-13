package app

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetApp(t *testing.T) {
	viper.SetDefault("logger.level", "DEBUG")
	viper.SetConfigFile("../config/test.yaml")
	app := GetApp("127.0.0.1", 9999, false)
	if app.Port != 9999 || app.Host != "127.0.0.1" {
		t.Fail()
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	viper.SetConfigFile("../config/invalid")
	app = GetApp("127.0.0.1", 9999, false)
}
