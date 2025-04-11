package errors

import "net/http"

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

// CodeMap is a map of go-service errors and http status codes.
var CodeMap = map[string]int{
	ECONFLICT:    http.StatusConflict,
	EINVALID:     http.StatusBadRequest,
	ENOTFOUND:    http.StatusNotFound,
	EINTERNAL:    http.StatusInternalServerError,
	EFORBIDDEN:   http.StatusForbidden,
	ENOTMODIFIED: http.StatusNotModified,
}
