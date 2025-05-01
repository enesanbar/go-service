package grpc

import (
	"context"
	"log"
	"net"
	"path/filepath"

	serviceLogger "github.com/enesanbar/go-service/log"
	"go.uber.org/zap"

	"google.golang.org/grpc/stats"
)

// StatsHandler implements [stats.StatsHandler](https://pkg.go.dev/google.golang.org/grpc/stats#StatsHandler) interface.
type StatsHandler struct {
	logger serviceLogger.Factory
}

type connStatCtxKey struct{}

// TagConn can attach some information to the given context.
// The context used in HandleConn for this connection will be derived from the context returned.
// In the gRPC client:
// The context used in HandleRPC for RPCs on this connection will be the user's context and NOT derived from the context returned here.
// In the gRPC server:
// The context used in HandleRPC for RPCs on this connection will be derived from the context returned here.
func (st *StatsHandler) TagConn(ctx context.Context, stat *stats.ConnTagInfo) context.Context {
	log.Printf("[TagConn] [%T]: %+[1]v", stat)
	return context.WithValue(ctx, connStatCtxKey{}, stat)
}

// HandleConn processes the Conn stats.
func (st *StatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	var rAddr net.Addr
	if s, ok := ctx.Value(connStatCtxKey{}).(*stats.ConnTagInfo); ok {
		rAddr = s.RemoteAddr
	}

	if stat.IsClient() {
		log.Printf("[server addr: %s] [HandleConn] [%T]: %+[2]v", rAddr, stat)
	} else {
		log.Printf("[client addr: %s] [HandleConn] [%T]: %+[2]v", rAddr, stat)
	}
}

type rpcStatCtxKey struct{}

// TagRPC can attach some information to the given context.
// The context used for the rest lifetime of the RPC will be derived from the returned context.
func (st *StatsHandler) TagRPC(ctx context.Context, stat *stats.RPCTagInfo) context.Context {
	log.Printf("[TagRPC] [%T]: %+[1]v", stat)
	return context.WithValue(ctx, rpcStatCtxKey{}, stat)
}

// HandleRPC processes the RPC stats. Note: All stat fields are read-only.
func (st *StatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	var sMethod string
	if s, ok := ctx.Value(rpcStatCtxKey{}).(*stats.RPCTagInfo); ok {
		sMethod = filepath.Base(s.FullMethodName)
	}

	var cAddr net.Addr
	// for gRPC clients, key connStatCtxKey{} will not be present in HandleRPC's context.
	if s, ok := ctx.Value(connStatCtxKey{}).(*stats.ConnTagInfo); ok {
		cAddr = s.RemoteAddr
	}

	// log the stats in json format
	st.logger.For(ctx).With(zap.Any("stats", stat)).Info("gRPC stats")
	if stat.IsClient() {
		log.Printf("[server method: %s] [HandleRPC] [%T]: %+[2]v", sMethod, stat)
	} else {
		log.Printf("[client addr: %s] [HandleRPC] [%T]: %+[2]v", cAddr, stat)
	}
}

// New returns a new implementation of [stats.Handler](https://pkg.go.dev/google.golang.org/grpc/stats#Handler) interface.
func NewStatsHandler(logger serviceLogger.Factory) *StatsHandler {
	return &StatsHandler{
		logger: logger,
	}
}
