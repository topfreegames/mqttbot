// mqttbot
// https://github.com/topfreegames/mqttbot
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package app

import (
	"fmt"

	"github.com/getsentry/raven-go"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/bot"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/topfreegames/mqttbot/redisclient"
	"github.com/uber-go/zap"
)

// App is the struct that defines the application
type App struct {
	Debug       bool
	Port        int
	Host        string
	Api         *iris.Framework
	MqttBot     *bot.MqttBot
	RedisClient *redisclient.RedisClient
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
	app.configureSentry()
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

func (app *App) configureSentry() {
	sentryURL := viper.GetString("sentry.url")
	logger.Logger.Info(fmt.Sprintf("Configuring sentry with URL %s", sentryURL))
	raven.SetDSN(sentryURL)
}

func (app *App) configureApplication() {
	app.MqttBot = bot.GetMqttBot()
	app.Api = iris.New()
	a := app.Api
	a.Use(NewLoggerMiddleware(zap.New(
		zap.NewJSONEncoder(),
	)))
	a.Use(&SentryMiddleware{App: app})
	a.Use(&VersionMiddleware{App: app})
	a.Use(&RecoveryMiddleware{OnError: app.onErrorHandler})
	a.Get("/healthcheck", HealthCheckHandler(app))
	a.Get("/history/*topic", HistoryHandler(app))

	app.RedisClient = redisclient.GetRedisClient(viper.GetString("redis.host"), viper.GetInt("redis.port"), viper.GetString("redis.password"))
}

func (app *App) onErrorHandler(err error, stack []byte) {
	logger.Logger.Errorf(
		"Panic occurred. stack: %s", string(stack),
	)
	tags := map[string]string{
		"source": "app",
		"type":   "panic",
	}
	raven.CaptureError(err, tags)
}

// Start starts the application
func (app *App) Start() {
	if viper.GetBool("api.tls") {
		logger.Logger.Infof("Api listening using TLS! Certfile: %s, Keyfile: %s", viper.GetString("api.certFile"), viper.GetString("api.keyFile"))
		app.Api.ListenTLS(fmt.Sprintf("%s:%d", app.Host, app.Port), viper.GetString("api.certFile"), viper.GetString("api.keyFile"))
	} else {
		app.Api.Listen(fmt.Sprintf("%s:%d", app.Host, app.Port))
	}
}
