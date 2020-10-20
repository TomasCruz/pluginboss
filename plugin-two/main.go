package main

import (
	"github.com/TomasCruz/pluginboss/plugininfo"

	"github.com/hashicorp/go-plugin"
)

var pluginName string = "conversion_two"

// Converter is an implementation of plugin's interface
type Converter struct{}

func (Converter) Convert(in float64) (out float64, err error) {
	// Celsius to Fahrenheit
	out = in*9/5 + 32
	return
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "PLUGIN_TWO",
			MagicCookieValue: "hello",
		},
		Plugins: map[string]plugin.Plugin{
			pluginName: &plugininfo.ConverterGRPCPlugin{Impl: &Converter{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
