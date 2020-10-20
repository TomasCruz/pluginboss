// plugin packages live under this repo for demonstration purposes only.
package main

import (
	"github.com/TomasCruz/pluginboss/plugininfo"

	"github.com/hashicorp/go-plugin"
)

var pluginName string = "conversion_one"

// Converter is an implementation of plugin's interface
type Converter struct{}

func (Converter) Convert(in float64) (out float64, err error) {
	// Fahrenheit to Celsius
	out = (in - 32) * 5 / 9
	return
}

func main() {
	// hardcoding handshake parameter. The idea is that plugin is provided by the 3rd party,
	// so it is the owner of handshake info any way. As the plugin also need to know about handshake, it's hardcoded here as well
	// Vision is having pluginboss running as HTTP server, loading plugin data on startup from a JSON,
	// having plugin data as provided by it's authors. It can however, potentially get or update plugin data via endpoints.
	// Having it in a file is not the finest solution, but it seems fine at the moment
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "PLUGIN_ONE",
			MagicCookieValue: "hello",
		},
		Plugins: map[string]plugin.Plugin{
			pluginName: &plugininfo.ConverterGRPCPlugin{Impl: &Converter{}},
		},
		GRPCServer: plugin.DefaultGRPCServer, // A non-nil value here enables gRPC serving for this plugin...
	})
}
