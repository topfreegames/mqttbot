package app

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/topfreegames/mqttbot/models"
	"github.com/topfreegames/mqttbot/plugins/modules"
	"github.com/yuin/gopher-lua"
	"os"
	"strings"
)

type Plugins struct {
	PluginMappings map[string]string
	Config         *viper.Viper
	LState         *lua.LState
}

func GetPlugins() *Plugins {
	plugins := &Plugins{
		Config:         viper.New(),
		PluginMappings: map[string]string{},
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
	p.PluginMappings = p.Config.GetStringMapString("plugins.mappings")
}

func (p *Plugins) loadModules() {
	L := p.loadPersistenceModule()
	if err := L.DoFile("plugins/load_modules.lua"); err != nil {
		logger.Logger.Error("Error loading lua go modules, err:", err)
		os.Exit(1)
	}
	p.LState = L
	logger.Logger.Info("Successfully loaded lua go modules")
}

func (p *Plugins) loadPersistenceModule() *lua.LState {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("persistence_module", modules.PersistenceModuleLoader)
	return L
}

func (p *Plugins) ExecutePluginWithMessage(plugin string, message *models.Message) {
	//todo checar se existe
}
