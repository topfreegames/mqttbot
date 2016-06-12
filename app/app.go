package app

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"strings"
)

type App struct {
	Debug      bool
	Port       int
	Host       string
	ConfigPath string
	Api        *iris.Iris
	MqttClient *MqttClient
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
	app.Config.SetDefault("mqttserver.host", "localhost")
	app.Config.SetDefault("mqttserver.port", 1883)
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
	app.MqttClient = GetMqttClient(app.Config.GetString("mqttserver.host"), app.Config.GetInt("mqttserver.port"))
	app.Api = iris.New()
	a := app.Api

	a.Get("/healthcheck", HealthCheckHandler(app))
}

func (app *App) Start() {
	app.MqttClient.Start()
	app.Api.Listen(fmt.Sprintf("%s:%d", app.Host, app.Port))
}
