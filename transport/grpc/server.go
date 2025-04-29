package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/enesanbar/go-service/log"
	"github.com/enesanbar/go-service/wiring"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

type GRPCServerParams struct {
	fx.In

	Logger log.Factory
	Config *ServerConfig
}

type GRPCServer struct {
	Server *grpc.Server
	logger log.Factory
	cfg    *ServerConfig
}

func New(p GRPCServerParams) (wiring.RunnableGroup, *GRPCServer) {
	s := grpc.NewServer()

	grpcServer := &GRPCServer{
		Server: s,
		logger: p.Logger,
		cfg:    p.Config,
	}

	return wiring.RunnableGroup{
		Runnable: grpcServer,
	}, grpcServer
}

func (s *GRPCServer) Start() error {
	s.logger.Bg().
		With(zap.Int("port", s.cfg.Port)).
		Info("starting GRPC Server")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		s.logger.Bg().With(zap.Error(err)).Error("failed to listen")
		return err
	}

	if err := s.Server.Serve(lis); err != nil {
		s.logger.Bg().With(zap.Error(err)).Error("failed to serve")
		return err
	}

	return nil
}

func (s *GRPCServer) Stop() error {
	timer := time.AfterFunc(time.Duration(s.cfg.GracefulStopTimeoutSeconds)*time.Second, func() {
		s.logger.Bg().Error("grpc server could not be stopped gracefully, forcing stop")
		s.Server.Stop()
	})
	defer timer.Stop()

	s.Server.GracefulStop()
	s.logger.Bg().Info("grpc server stopped gracefully")
	return nil
}
