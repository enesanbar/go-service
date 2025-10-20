package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

// ServerParams holds the parameters for creating a gRPC server.
type ServerParams struct {
	fx.In

	Logger log.Factory
	Config *ServerConfig

	ServerOptions []grpc.ServerOption `group:"grpc-server-options"`
}

// Server is a wrapper around the gRPC server.
type Server struct {
	Server *grpc.Server
	logger log.Factory
	cfg    *ServerConfig
}

// NewServer creates a new gRPC server with the provided parameters.
func NewServer(p ServerParams) (*Server, error) {
	s := grpc.NewServer(p.ServerOptions...)
	reflection.Register(s)
	grpcServer := &Server{
		Server: s,
		logger: p.Logger,
		cfg:    p.Config,
	}

	return grpcServer, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.For(ctx).
		With(zap.Int("port", s.cfg.Port)).
		With(zap.Any("config", s.cfg)).
		Info("starting GRPC Server")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		s.logger.For(ctx).With(zap.Error(err)).Error("failed to listen")
		return err
	}

	if err := s.Server.Serve(lis); err != nil {
		s.logger.For(ctx).With(zap.Error(err)).Error("failed to serve")
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	timer := time.AfterFunc(time.Duration(s.cfg.GracefulStopTimeoutSeconds)*time.Second, func() {
		s.logger.Bg().Info("gRPC server could not be stopped gracefully, forcing stop")
		s.Server.Stop()
		s.logger.Bg().Info("http server forced to stop")
	})
	defer timer.Stop()

	s.logger.Bg().Info("gracefully stopping gRPC Server")
	s.Server.GracefulStop()
	s.logger.Bg().Info("gRPC server stopped gracefully")
	return nil
}
