package app

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/bot"
	"github.com/topfreegames/mqttbot/logger"
)

// App is the struct that defines the application
type App struct {
	Debug   bool
	Port    int
	Host    string
	Api     *iris.Framework
	MqttBot *bot.MqttBot
}

// GetApp creates an app given the parameters
func GetApp(host string, port int, debug bool) *App {
	logger.SetupLogger(viper.GetString("logger.level"))
	logger.Logger.Debug(
		fmt.Sprintf("Starting app with host: %s, port: %d", host, port))
	app := &App{
		Host:  host,
		Port:  port,
		Debug: debug,
	}
	app.Configure()
	return app
}

// Configure configures the application
func (app *App) Configure() {
	app.setConfigurationDefaults()
	app.loadConfiguration()
	app.configureApplication()
}

func (app *App) setConfigurationDefaults() {
	viper.SetDefault("healthcheck.workingText", "WORKING")
}

func (app *App) loadConfiguration() {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logger.Logger.Debug("Config file read successfully")
	} else {
		panic(fmt.Sprintf("Could not load configuration file"))
	}
}

func (app *App) configureApplication() {
	app.MqttBot = bot.GetMqttBot()
	app.Api = iris.New()
	a := app.Api

	a.Get("/healthcheck", HealthCheckHandler(app))
}

// Start starts the application
func (app *App) Start() {
	app.Api.Listen(fmt.Sprintf("%s:%d", app.Host, app.Port))
}
