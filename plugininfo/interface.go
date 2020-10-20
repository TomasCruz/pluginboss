// Package plugininfo contains type definitions and some boilerplating Go plugins can use
// The idea is to have it in a separate repo as a library, used by both pluginboss and (Go) plugins
// Plugins in other languages need to implement it
package plugininfo

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
)

// ConverterPlugin is the plugin's interface. Name is deliberatelly wrong (ConverterPlugin instead of Converter) to
// emphasize the fact it's an interface defining plugin API
type ConverterPlugin interface {
	Convert(in float64) (float64, error)
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type ConverterGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins that are written in Go.
	Impl ConverterPlugin
}

func (p *ConverterGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterConverterServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *ConverterGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: NewConverterClient(c)}, nil
}
