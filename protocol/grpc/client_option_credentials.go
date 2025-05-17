package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientOptionCredentialsParams struct {
	fx.In

	Logger log.Factory
	Config *ServerConfig
}

func NewClientOptionCredentials(p ClientOptionCredentialsParams) grpc.DialOption {
	if !p.Config.TLS.ClientTLSEnabled {
		return grpc.WithTransportCredentials(insecure.NewCredentials())
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
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caPool,
		InsecureSkipVerify: false,
	}
	creds := credentials.NewTLS(tlsCfg)

	return grpc.WithTransportCredentials(creds)
}
