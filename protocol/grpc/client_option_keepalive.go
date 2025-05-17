package grpc

import (
	"time"

	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// ClientOptionKeepaliveParams holds the parameters for creating a gRPC dial option for keepalive parameters.
type ClientOptionKeepaliveParams struct {
	fx.In

	Logger log.Factory
	Config *ServerConfig
}

// NewClientOptionKeepAliveParams creates a new gRPC dial option for keepalive parameters.
func NewClientOptionKeepAliveParams(p ClientOptionKeepaliveParams) grpc.DialOption {
	var kasp = keepalive.ClientParameters{
		Time:                time.Duration(p.Config.KeepAlive.TimeSeconds) * time.Second,
		Timeout:             time.Duration(p.Config.KeepAlive.TimeoutSeconds) * time.Second,
		PermitWithoutStream: p.Config.KeepAlive.PermitWithoutStream,
	}
	return grpc.WithKeepaliveParams(kasp)
}
