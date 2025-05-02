package grpc

import (
	"time"

	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type GRPCClientOptionKeepaliveParams struct {
	fx.In

	Logger log.Factory
	Config *ServerConfig
}

func NewGRPCClientOptionKeepAliveParams(p GRPCClientOptionKeepaliveParams) grpc.DialOption {
	var kasp = keepalive.ClientParameters{
		Time:                time.Duration(p.Config.KeepAlive.TimeSeconds) * time.Second,
		Timeout:             time.Duration(p.Config.KeepAlive.TimeoutSeconds) * time.Second,
		PermitWithoutStream: p.Config.KeepAlive.PermitWithoutStream,
	}
	return grpc.WithKeepaliveParams(kasp)
}
