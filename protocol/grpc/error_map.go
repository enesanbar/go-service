package grpc

import (
	coreErr "github.com/enesanbar/go-service/core/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CodeMapGRPC is a map of go-service errors and grpc status codes.
var CodeMapGRPC = map[string]codes.Code{
	coreErr.ECONFLICT:    codes.FailedPrecondition,
	coreErr.EINVALID:     codes.InvalidArgument,
	coreErr.ENOTFOUND:    codes.NotFound,
	coreErr.EINTERNAL:    codes.Internal,
	coreErr.EFORBIDDEN:   codes.PermissionDenied,
	coreErr.ENOTMODIFIED: codes.OK,
}

// ErrorStatus recursively checks err.Code and
// returns appropriate grpc response code depending on the error,
// otherwise it returns 500
func ErrorStatus(err error) codes.Code {
	if err == nil {
		return codes.Internal
	} else if e, ok := err.(*coreErr.Error); ok && e.Code != "" {
		return CodeMapGRPC[e.Code]
	} else if ok && e.Err != nil {
		return ErrorStatus(e.Err)
	} else if statusErr, ok := status.FromError(err); ok {
		return statusErr.Code()
	}
	return codes.Internal
}
