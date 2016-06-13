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

func PersistenceModuleLoader(L *lua.LState) int {
	configurePersistenceModule()
	configureDatabase()
	mod := L.SetFuncs(L.NewTable(), persistenceModuleExports)
	L.Push(mod)
	return 1
}

var persistenceModuleExports = map[string]lua.LGFunction{
	"index_message":  IndexMessage,
	"query_messages": QueryMessages,
}

func configurePersistenceModule() {
	setConfigurationDefaults()
}

func setConfigurationDefaults() {
	viper.SetDefault("elasticsearch.host", "localhost")
	viper.SetDefault("elasticsearch.port", 9200)
	viper.SetDefault("elasticsearch.sniff", false)
	viper.SetDefault("elasticsearch.indexMappings", map[string]string{})
}

func configureDatabase() {
	logger.Logger.Debug(fmt.Sprintf("Connecting to elasticsearch @ http://%s:%d", viper.GetString("elasticsearch.host"), viper.GetInt("elasticsearch.port")))
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%d", viper.GetString("elasticsearch.host"), viper.GetInt("elasticsearch.port"))),
		elastic.SetSniff(viper.GetBool("elasticsearch.sniff")),
	)
	if err != nil {
		logger.Logger.Fatal("Failed to connect to elasticsearch! err:", err)
	}
	logger.Logger.Info(fmt.Sprintf("Successfully connected to elasticsearch @ http://%s:%d", viper.GetString("elasticsearch.host"), viper.GetInt("elasticsearch.port")))
	logger.Logger.Debug("Creating index chat into ES")

	indexes := viper.GetStringMapString("elasticsearch.indexMappings")
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
