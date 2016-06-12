package app

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/spf13/viper"
	"gopkg.in/olivere/elastic.v3"
	"os"
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
	ESClient   *elastic.Client
}

func GetApp(host string, port int, configPath string, debug bool) *App {
	SetupLogger()
	Logger.Debug(fmt.Sprintf("Starting app with host: %s, port: %d, configFile: %s", host, port, configPath))
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
	app.configureDatabase()
	app.configureApplication()
}

func (app *App) setConfigurationDefaults() {
	app.Config.SetDefault("healthcheck.workingText", "WORKING")
	app.Config.SetDefault("elasticsearch.host", "localhost")
	app.Config.SetDefault("elasticsearch.port", 9200)
	app.Config.SetDefault("elasticsearch.sniff", false)
	app.Config.SetDefault("mqttserver.host", "localhost")
	app.Config.SetDefault("mqttserver.port", 1883)
}

func (app *App) loadConfiguration() {
	app.Config.SetConfigFile(app.ConfigPath)
	app.Config.SetEnvPrefix("mqttbot")
	app.Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	app.Config.AutomaticEnv()

	if err := app.Config.ReadInConfig(); err == nil {
		Logger.Debug(fmt.Sprintf("Using config file: %s", app.Config.ConfigFileUsed()))
	}
}

func (app *App) configureDatabase() {
	Logger.Debug(fmt.Sprintf("Connecting to elasticsearch @ http://%s:%d", app.Config.GetString("elasticsearch.host"), app.Config.GetInt("elasticsearch.port")))
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%d", app.Config.GetString("elasticsearch.host"), app.Config.GetInt("elasticsearch.port"))),
		elastic.SetSniff(app.Config.GetBool("elasticsearch.sniff")),
	)
	if err != nil {
		Logger.Error("Failed to connect to elasticsearch! err:", err)
		os.Exit(1)
	}
	Logger.Info(fmt.Sprintf("Successfully connected to elasticsearch @ http://%s:%d", app.Config.GetString("elasticsearch.host"), app.Config.GetInt("elasticsearch.port")))
	Logger.Debug("Creating index chat into ES")

	indexMapping := `
		{
			"mappings": {
				"message": {
					"_timestamp": { 
						"enabled": true
					},
					"_ttl": { 
						"enabled": true,
						"default": "2d"
					}
				}
			}
		}	
	`
	_, err = client.CreateIndex("chat").Body(indexMapping).Do()
	if err != nil {
		if strings.Contains(err.Error(), "index_already_exists_exception") {
			Logger.Warning("Chat index already exists into ES! Ignoring creation...")
		} else {
			Logger.Error("Failed to create chat index into ES, err:", err)
			os.Exit(1)
		}
	} else {
		Logger.Debug("Sucessfully created index chat into ES")
	}
	app.ESClient = client
}

func (app *App) configureApplication() {
	app.MqttClient = GetMqttClient(app.Config.GetString("mqttserver.host"), app.Config.GetInt("mqttserver.port"), app.ESClient)
	app.Api = iris.New()
	a := app.Api

	if app.Debug {
		a.Use(logger.New(iris.Logger()))
	}

	a.Get("/healthcheck", HealthCheckHandler(app))
}

func (app *App) Start() {
	app.MqttClient.Start()
	app.Api.Listen(fmt.Sprintf("%s:%d", app.Host, app.Port))
}
