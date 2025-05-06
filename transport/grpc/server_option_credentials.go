package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCServerOptionCredentialsParams struct {
	fx.In

	Logger log.Factory
	Config *ServerConfig
}

func NewGRPCServerOptionCredentials(p GRPCServerOptionCredentialsParams) grpc.ServerOption {
	if !p.Config.TLS.ServerTLSEnabled {
		return grpc.Creds(insecure.NewCredentials())
	}

	cert, err := tls.LoadX509KeyPair(p.Config.TLS.CertFile, p.Config.TLS.KeyFile)
	if err != nil {
		p.Logger.Bg().With(zap.Error(err)).Fatal("failed to load TLS cert and key")
	}

	caCert, err := os.ReadFile(p.Config.TLS.CAFile)
	if err != nil {
		p.Logger.Bg().With(zap.Error(err)).Fatal("failed to read CA cert")
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		p.Logger.Bg().With(zap.Error(err)).Fatal("failed to append CA cert")
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	}
	creds := credentials.NewTLS(tlsCfg)

	return grpc.Creds(creds)
}
