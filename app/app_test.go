package app

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetApp(t *testing.T) {
	viper.SetDefault("logger.level", "DEBUG")
	viper.SetConfigFile("../config/test.yml")
	app := GetApp("127.0.0.1", 9999, "../config/test.yml", false)
	if app.Port != 9999 || app.Host != "127.0.0.1" {
		t.Fail()
	}
}
