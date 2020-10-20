package main

import (
	"github.com/TomasCruz/pluginboss/plugininfo"

	"github.com/hashicorp/go-plugin"
)

var pluginName string = "conversion_three"

// Converter is an implementation of plugin's interface
type Converter struct{}

func (Converter) Convert(in float64) (out float64, err error) {
	panic("blah blah") // whaddya mean Convert
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "PLUGIN_THREE",
			MagicCookieValue: "hello",
		},
		Plugins: map[string]plugin.Plugin{
			pluginName: &plugininfo.ConverterGRPCPlugin{Impl: &Converter{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
