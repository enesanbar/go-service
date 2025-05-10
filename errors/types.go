package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// transport-agnostic error types, cli, rest, gprc, graphql.
// at least one code should be provided in the error in the chain of errors
const (
	ECONFLICT    = "conflict"     // action cannot be performed
	EINVALID     = "invalid"      // validation failed
	ENOTFOUND    = "not_found"    // entity does not exist
	EINTERNAL    = "internal"     // internal error
	EFORBIDDEN   = "unauthorized" // unauthorized
	ENOTMODIFIED = "not_modified" // unauthorized
)

// CodeMapHTTP is a map of go-service errors and http status codes.
var CodeMapHTTP = map[string]int{
	ECONFLICT:    http.StatusConflict,
	EINVALID:     http.StatusBadRequest,
	ENOTFOUND:    http.StatusNotFound,
	EINTERNAL:    http.StatusInternalServerError,
	EFORBIDDEN:   http.StatusForbidden,
	ENOTMODIFIED: http.StatusNotModified,
}

// CodeMapGRPC is a map of go-service errors and grpc status codes.
var CodeMapGRPC = map[string]codes.Code{
	ECONFLICT:    codes.FailedPrecondition,
	EINVALID:     codes.InvalidArgument,
	ENOTFOUND:    codes.NotFound,
	EINTERNAL:    codes.Internal,
	EFORBIDDEN:   codes.PermissionDenied,
	ENOTMODIFIED: codes.OK,
}
