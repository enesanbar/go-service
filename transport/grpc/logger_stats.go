package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/enesanbar/go-service/info"
	"github.com/enesanbar/go-service/log"
	"go.uber.org/zap"

	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// LoggerStatsHandler implements [stats.LoggerStatsHandler](https://pkg.go.dev/google.golang.org/grpc/stats#LoggerStatsHandler) interface.
type LoggerStatsHandler struct {
	logger log.Factory
}

func NewStatsHandler(logger log.Factory) *LoggerStatsHandler {
	return &LoggerStatsHandler{
		logger: logger,
	}
}

type connStatCtxKey struct{}

func (st *LoggerStatsHandler) TagConn(ctx context.Context, stat *stats.ConnTagInfo) context.Context {
	return context.WithValue(ctx, connStatCtxKey{}, stat)
}

func (st *LoggerStatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
}

type rpcStatCtxKey struct{}

func (st *LoggerStatsHandler) TagRPC(ctx context.Context, stat *stats.RPCTagInfo) context.Context {
	return context.WithValue(ctx, rpcStatCtxKey{}, stat)
}

func (st *LoggerStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	var sMethod string
	if s, ok := ctx.Value(rpcStatCtxKey{}).(*stats.RPCTagInfo); ok {
		sMethod = filepath.Base(s.FullMethodName)
	}

	logPayload := func(fieldName string, payload interface{}) {
		var payloadJSON string

		if msg, ok := payload.(proto.Message); ok {
			b, err := protojson.Marshal(msg)
			if err == nil {
				payloadJSON = string(b)
			} else {
				payloadJSON = err.Error()
			}
		} else {
			b, err := json.Marshal(payload)
			if err == nil {
				payloadJSON = string(b)
			} else {
				payloadJSON = err.Error()
			}
		}

		st.logger.For(ctx).Info(
			fmt.Sprintf("gRPC %s payload", fieldName),
			zap.String("service", info.ServiceName),
			zap.String("method", sMethod),
			zap.String(fieldName, payloadJSON),
		)
	}

	switch s := stat.(type) {
	case *stats.InPayload:
		logPayload("request", s.Payload)
	case *stats.OutPayload:
		logPayload("response", s.Payload)
	case *stats.End:
		fields := []zap.Field{
			zap.String("service", sMethod),
			zap.String("method", sMethod),
		}

		// Include status code and error message
		if s.Error != nil {
			sts, _ := status.FromError(s.Error)
			fields = append(fields,
				zap.String("status", sts.Code().String()),
				zap.String("error", sts.Message()),
			)
		} else {
			fields = append(fields, zap.String("status", "OK"))
		}

		st.logger.For(ctx).Info("gRPC call finished", fields...)
	}
}
