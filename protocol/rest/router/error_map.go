package router

import (
	"net/http"

	"github.com/enesanbar/go-service/core/errors"
)

// CodeMapHTTP is a map of go-service errors and http status codes.
var CodeMapHTTP = map[string]int{
	errors.ECONFLICT:    http.StatusConflict,
	errors.EINVALID:     http.StatusBadRequest,
	errors.ENOTFOUND:    http.StatusNotFound,
	errors.EINTERNAL:    http.StatusInternalServerError,
	errors.EFORBIDDEN:   http.StatusForbidden,
	errors.ENOTMODIFIED: http.StatusNotModified,
}

// ErrorStatus recursively checks err.Code and
// returns appropriate http response code depending on the error,
// otherwise it returns 500
func ErrorStatus(err error) int {
	if err == nil {
		return 500
	} else if e, ok := err.(*errors.Error); ok && e.Code != "" {
		return CodeMapHTTP[e.Code]
	} else if ok && e.Err != nil {
		return ErrorStatus(e.Err)
	}
	return 500
}
