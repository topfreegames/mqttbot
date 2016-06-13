package app

import (
	"fmt"
	"strings"

	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/topfreegames/mqttbot/mqtt"
)

type App struct {
	Debug      bool
	Port       int
	Host       string
	ConfigPath string
	Api        *iris.Iris
	MqttClient *mqtt.MqttClient
	Config     *viper.Viper
}

func GetApp(host string, port int, configPath string, debug bool) *App {
	logger.SetupLogger()
	logger.Logger.Debug(fmt.Sprintf("Starting app with host: %s, port: %d, configFile: %s", host, port, configPath))
	app := &App{
		Host:       host,
		Port:       port,
		ConfigPath: configPath,
		Config:     viper.New(),
		Debug:      debug,
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
	app.Config.SetDefault("healthcheck.workingText", "WORKING")
}

func (app *App) loadConfiguration() {
	app.Config.SetConfigFile(app.ConfigPath)
	app.Config.SetEnvPrefix("mqttbot")
	app.Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	app.Config.AutomaticEnv()

	if err := app.Config.ReadInConfig(); err == nil {
		logger.Logger.Debug(fmt.Sprintf("Using config file: %s", app.Config.ConfigFileUsed()))
	}
}

func (app *App) configureApplication() {
	app.MqttClient = mqtt.GetMqttClient()
	app.Api = iris.New()
	a := app.Api

	a.Get("/healthcheck", HealthCheckHandler(app))
}

func (app *App) Start() {
	app.Api.Listen(fmt.Sprintf("%s:%d", app.Host, app.Port))
}
