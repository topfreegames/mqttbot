package es

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"gopkg.in/topfreegames/elastic.v2"
)

var esclient *elastic.Client
var once sync.Once

// GetESClient returns the elasticsearch client with the given configs
func GetESClient() *elastic.Client {
	once.Do(func() {
		configure()
	})
	return esclient
}

func configure() {
	setConfigurationDefaults()
	configureESClient()
}

func setConfigurationDefaults() {
	viper.SetDefault("elasticsearch.host", "http://localhost:9200")
	viper.SetDefault("elasticsearch.sniff", false)
	viper.SetDefault("elasticsearch.indexMappings", map[string]string{})
}

func configureESClient() {
	logger.Logger.Debug(fmt.Sprintf("Connecting to elasticsearch @ %s",
		viper.GetString("elasticsearch.host")))
	credentials := defaults.CredChain(defaults.Config(), defaults.Handlers())
	awsSigningRoundTripper := elastic.NewAWSSigningRoundTripper(nil, "us-east-1", credentials)
	esHTTPClient := &http.Client{Transport: awsSigningRoundTripper}

	client, err := elastic.NewClient(
		elastic.SetHttpClient(esHTTPClient),
		elastic.SetURL(viper.GetString("elasticsearch.host")),
		elastic.SetSniff(viper.GetBool("elasticsearch.sniff")),
	)
	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Failed to connect to elasticsearch! err: %v", err))
	}
	logger.Logger.Info(fmt.Sprintf("Successfully connected to elasticsearch @ %s",
		viper.GetString("elasticsearch.host")))
	logger.Logger.Debug("Creating index chat into ES")

	indexes := viper.GetStringMapString("elasticsearch.indexMappings")
	for index, mappings := range indexes {
		_, err = client.CreateIndex(index).Body(mappings).Do()
		if err != nil {
			if strings.Contains(err.Error(), "index_already_exists_exception") || strings.Contains(err.Error(), "IndexAlreadyExistsException") {
				logger.Logger.Warning(fmt.Sprintf("Index %s already exists into ES! Ignoring creation...", index))
			} else {
				logger.Logger.Error(fmt.Sprintf("Failed to create index %s into ES, err: %s", index, err))
				os.Exit(1)
			}
		} else {
			logger.Logger.Debug(fmt.Sprintf("Sucessfully created index %s into ES", index))
		}
	}
	esclient = client
}
