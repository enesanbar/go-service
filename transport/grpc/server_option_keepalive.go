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
	// TODO: get these values from config. p.Config
	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	return grpc.KeepaliveEnforcementPolicy(kaep)
}

func NewGRPCServerOptionKeepAliveParams(p GRPCServerOptionKeepaliveParams) grpc.ServerOption {
	// TODO: get these values from config. p.Config
	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}
	return grpc.KeepaliveParams(kasp)
}
