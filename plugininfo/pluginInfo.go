package plugininfo

import (
	"github.com/hashicorp/go-plugin"
)

type PluginData struct {
	Cmd       string
	Handshake plugin.HandshakeConfig
}

type PluginInfo struct {
	PluginMap     map[string]plugin.Plugin
	PluginDataMap map[string]PluginData
}
