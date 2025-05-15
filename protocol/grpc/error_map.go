package grpc

import (
	"github.com/enesanbar/go-service/core/errors"
	"google.golang.org/grpc/codes"
)

// CodeMapGRPC is a map of go-service errors and grpc status codes.
var CodeMapGRPC = map[string]codes.Code{
	errors.ECONFLICT:    codes.FailedPrecondition,
	errors.EINVALID:     codes.InvalidArgument,
	errors.ENOTFOUND:    codes.NotFound,
	errors.EINTERNAL:    codes.Internal,
	errors.EFORBIDDEN:   codes.PermissionDenied,
	errors.ENOTMODIFIED: codes.OK,
}

// ErrorStatusGRPC recursively checks err.Code and
// returns appropriate grpc response code depending on the error,
// otherwise it returns 500
func ErrorStatus(err error) codes.Code {
	if err == nil {
		return 500
	} else if e, ok := err.(errors.Error); ok && e.Code != "" {
		return CodeMapGRPC[e.Code]
	} else if ok && e.Err != nil {
		return ErrorStatus(e.Err)
	}
	return 500
}
