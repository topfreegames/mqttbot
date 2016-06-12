package app

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/topfreegames/mqttbot/models"
	"github.com/topfreegames/mqttbot/plugins/modules"
	"github.com/yuin/gopher-lua"
)

type Plugins struct {
	PluginMappings []map[string]string
	Config         *viper.Viper
	LState         *lua.LState
}

func GetPlugins() *Plugins {
	plugins := &Plugins{
		Config:         viper.New(),
		PluginMappings: []map[string]string{},
	}
	return plugins
}

func (p *Plugins) SetupPlugins() {
	p.setConfigurationDefaults()
	p.loadConfiguration()
	p.loadPluginsMappings()
	p.loadModules()
}

func (p *Plugins) setConfigurationDefaults() {
	p.Config.SetDefault("plugins.mappings", map[string]string{})
}

func (p *Plugins) loadConfiguration() {
	p.Config.SetConfigFile("./config/plugins.yaml")
	p.Config.SetEnvPrefix("mqttbot.plugins")
	p.Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	p.Config.AutomaticEnv()

	if err := p.Config.ReadInConfig(); err == nil {
		logger.Logger.Debug(fmt.Sprintf("Using config file: %s", p.Config.ConfigFileUsed()))
	}
}

func (p *Plugins) loadPluginsMappings() {
	m := p.Config.GetStringMap("plugins")["mappings"]
	switch reflect.TypeOf(m).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(m)
		for i := 0; i < s.Len(); i++ {
			mappingInterface := s.Index(i).Interface()
			mappingAux := mappingInterface.(map[interface{}]interface{})
			mapping := make(map[string]string)
			mapping["topic"] = mappingAux[string("topic")].(string)
			mapping["messagePattern"] = mappingAux[string("messagePattern")].(string)
			mapping["pluginPath"] = mappingAux[string("pluginPath")].(string)
			p.PluginMappings = append(p.PluginMappings, mapping)
		}
	}
	logger.Logger.Debug(fmt.Sprintf("Loaded plugins mappings: %s", p.PluginMappings))

}

func (p *Plugins) loadModules() {
	L := p.preloadPersistenceModule()
	if err := L.DoFile("plugins/load_modules.lua"); err != nil {
		logger.Logger.Error("Error loading lua go modules, err:", err)
		os.Exit(1)
	}
	p.LState = L
	logger.Logger.Info("Successfully loaded lua go modules")
}

func (p *Plugins) preloadPersistenceModule() *lua.LState {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("persistence_module", modules.PersistenceModuleLoader)
	return L
}

func (p *Plugins) ExecutePluginWithMessage(plugin string, message *models.Message) {
	//todo checar se existe
}
