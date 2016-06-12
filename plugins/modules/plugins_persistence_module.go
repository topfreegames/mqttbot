package modules

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/yuin/gopher-lua"
	"gopkg.in/olivere/elastic.v3"
)

var esClient *elastic.Client
var config *viper.Viper

func PersistenceModuleLoader(L *lua.LState) int {
	config = viper.New()
	configurePersistenceModule()
	configureDatabase()
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"index_message":  IndexMessage,
	"query_messages": QueryMessages,
}

func configurePersistenceModule() {
	setConfigurationDefaults()
	loadConfiguration()
}

func setConfigurationDefaults() {
	config.SetDefault("es.host", "localhost")
	config.SetDefault("es.port", 9200)
	config.SetDefault("es.sniff", false)
	config.SetDefault("es.indexMappings", map[string]string{})
}

func loadConfiguration() {
	config.SetConfigFile("./config/elasticsearch.yaml")
	config.SetEnvPrefix("mqttbot.persistence")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err == nil {
		logger.Logger.Debug(fmt.Sprintf("Using config file: %s", config.ConfigFileUsed()))
	}
}

func configureDatabase() {
	logger.Logger.Debug(fmt.Sprintf("Connecting to elasticsearch @ http://%s:%d", config.GetString("es.host"), config.GetInt("es.port")))
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%d", config.GetString("es.host"), config.GetInt("es.port"))),
		elastic.SetSniff(config.GetBool("es.sniff")),
	)
	if err != nil {
		logger.Logger.Error("Failed to connect to elasticsearch! err:", err)
		os.Exit(1)
	}
	logger.Logger.Info(fmt.Sprintf("Successfully connected to elasticsearch @ http://%s:%d", config.GetString("es.host"), config.GetInt("es.port")))
	logger.Logger.Debug("Creating index chat into ES")

	indexes := config.GetStringMapString("es.indexMappings")
	for index, mappings := range indexes {
		_, err = client.CreateIndex(index).Body(mappings).Do()
		if err != nil {
			if strings.Contains(err.Error(), "index_already_exists_exception") {
				logger.Logger.Warning(fmt.Sprintf("Index %s already exists into ES! Ignoring creation...", index))
			} else {
				logger.Logger.Error(fmt.Sprintf("Failed to create index %s into ES, err: %s", index, err))
				os.Exit(1)
			}
		} else {
			logger.Logger.Debug(fmt.Sprintf("Sucessfully created index %s into ES", index))
		}
	}
	esClient = client
}

func IndexMessage(L *lua.LState) int {
	return 0
}

func QueryMessages(L *lua.LState) int {
	return 0
}
