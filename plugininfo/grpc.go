package plugininfo

import (
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of ConverterPlugin that talks over RPC.
// Basically a client with RPC communication details abstracted away
type GRPCClient struct {
	client ConverterClient
}

func (m *GRPCClient) Convert(in float64) (out float64, err error) {
	cresp, err := m.client.Convert(context.Background(), &ConvertRequest{In: in})
	if err != nil {
		return
	}

	out = cresp.Out
	return
}

// GRPCServer is the RPC server that GRPCClient talks to, a wrapper around plugin's server
type GRPCServer struct {
	Impl ConverterPlugin // plugin implementation
}

func (m *GRPCServer) Convert(
	ctx context.Context,
	req *ConvertRequest) (cresp *ConvertResponse, err error) {

	var out float64
	out, err = m.Impl.Convert(req.In)
	if err != nil {
		return
	}

	cresp = &ConvertResponse{Out: out}
	return
}
