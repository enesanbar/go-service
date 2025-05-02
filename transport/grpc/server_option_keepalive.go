package grpc

import (
	"time"

	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type GRPCServerOptionKeepaliveParams struct {
	fx.In

	Logger log.Factory
	Config *ServerConfig
}

func NewGRPCServerOptionKeepAliveEnforcementPolicy(p GRPCServerOptionKeepaliveParams) grpc.ServerOption {
	var kaep = keepalive.EnforcementPolicy{
		MinTime:             time.Duration(p.Config.KeepAlive.MinTimeSeconds) * time.Second,
		PermitWithoutStream: p.Config.KeepAlive.PermitWithoutStream,
	}

	return grpc.KeepaliveEnforcementPolicy(kaep)
}

func NewGRPCServerOptionKeepAliveParams(p GRPCServerOptionKeepaliveParams) grpc.ServerOption {
	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(p.Config.KeepAlive.MaxConnectionIdleSeconds) * time.Second,
		MaxConnectionAge:      time.Duration(p.Config.KeepAlive.MaxConnectionAgeSeconds) * time.Second,
		MaxConnectionAgeGrace: time.Duration(p.Config.KeepAlive.MaxConnectionAgeGraceSeconds) * time.Second,
		Time:                  time.Duration(p.Config.KeepAlive.TimeSeconds) * time.Second,
		Timeout:               time.Duration(p.Config.KeepAlive.TimeoutSeconds) * time.Second,
	}
	return grpc.KeepaliveParams(kasp)
}
