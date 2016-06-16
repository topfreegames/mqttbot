package app

import (
	"fmt"
	"strings"

	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/topfreegames/mqttbot/mqtt/bot"
)

type App struct {
	Debug   bool
	Port    int
	Host    string
	Api     *iris.Framework
	MqttBot *bot.MqttBot
}

func GetApp(host string, port int, debug bool) *App {
	logger.SetupLogger()
	logger.Logger.Debug(fmt.Sprintf("Starting app with host: %s, port: %d, configFile: %s", host, port, viper.ConfigFileUsed()))
	app := &App{
		Host:  host,
		Port:  port,
		Debug: debug,
	}
	app.Configure()
	return app
}

func (app *App) Configure() {
	app.setConfigurationDefaults()
	app.loadConfiguration()
	app.configureApplication()
}

func (app *App) setConfigurationDefaults() {
	viper.SetDefault("healthcheck.workingText", "WORKING")
}

func (app *App) loadConfiguration() {
	viper.SetEnvPrefix("mqttbot")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logger.Logger.Debug(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}
}

func (app *App) configureApplication() {
	app.MqttBot = bot.GetMqttBot()
	app.Api = iris.New()
	a := app.Api

	a.Get("/healthcheck", HealthCheckHandler(app))
}

func (app *App) Start() {
	app.Api.Listen(fmt.Sprintf("%s:%d", app.Host, app.Port))
}
