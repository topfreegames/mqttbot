// mqttbot
// https://github.com/topfreegames/mqttbot
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package app

import (
	"fmt"
	"os"

	"github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
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
	API         *echo.Echo
	Engine      engine.Server
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
		panic(fmt.Sprintf("Could not load configuration file, err: %s", err))
	}
}

func (app *App) configureSentry() {
	sentryURL := viper.GetString("sentry.url")
	logger.Logger.Info(fmt.Sprintf("Configuring sentry with URL %s", sentryURL))
	raven.SetDSN(sentryURL)
}

func (app *App) configureApplication() {
	app.MqttBot = bot.GetMqttBot()
	app.Engine = standard.New(fmt.Sprintf("%s:%d", app.Host, app.Port))
	app.API = echo.New()
	a := app.API
	_, w, _ := os.Pipe()
	a.SetLogOutput(w)
	a.Use(NewLoggerMiddleware(zap.New(
		zap.NewJSONEncoder(),
	)).Serve)
	a.Use(NewSentryMiddleware(app).Serve)
	a.Use(VersionMiddleware)
	a.Use(NewRecoveryMiddleware(app.OnErrorHandler).Serve)
	a.Get("/healthcheck", HealthCheckHandler(app))
	a.Get("/historysince/*", HistorySinceHandler(app))
	a.Get("/history/*", HistoryHandler(app))
	a.Get("/:other", NotFoundHandler(app))

	app.RedisClient = redisclient.GetRedisClient(viper.GetString("redis.host"), viper.GetInt("redis.port"), viper.GetString("redis.password"))
}

//OnErrorHandler handles application panics
func (app *App) OnErrorHandler(err interface{}, stack []byte) {
	logger.Logger.Error(err)

	var e error
	switch err.(type) {
	case error:
		e = err.(error)
	default:
		e = fmt.Errorf("%v", err)
	}

	tags := map[string]string{
		"source": "app",
		"type":   "panic",
	}
	raven.CaptureError(e, tags)
}

// Start starts the application
func (app *App) Start() {
	app.API.Run(app.Engine)
}
