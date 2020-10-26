package plugininfo

import (
	"github.com/hashicorp/go-plugin"
)

type PluginData struct {
	Cmd       string
	Handshake plugin.HandshakeConfig
	Plugin    plugin.Plugin
}

type PluginInfo struct {
	PluginDataMap map[string]PluginData
}
