package plugins

import (
	"github.com/topfreegames/mqttbot/logger"
	"github.com/topfreegames/mqttbot/plugins/modules"
	"github.com/yuin/gopher-lua"
)

type Plugins struct {
	PluginMappings []map[string]string
	LState         *lua.LState
}

func GetPlugins() *Plugins {
	plugins := &Plugins{
		PluginMappings: []map[string]string{},
	}
	return plugins
}

func (p *Plugins) SetupPlugins() {
	p.loadModules()
}

func (p *Plugins) loadModules() {
	L := p.preloadModules()
	if err := L.DoFile("plugins/load_modules.lua"); err != nil {
		logger.Logger.Fatal("Error loading lua go modules, err:", err)
	}
	p.LState = L
	logger.Logger.Info("Successfully loaded lua go modules")
}

func (p *Plugins) preloadModules() *lua.LState {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("persistence_module", modules.PersistenceModuleLoader)
	L.PreloadModule("mqttclient_module", modules.MqttClientModuleLoader)
	return L
}
