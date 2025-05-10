package errors

import "google.golang.org/grpc/codes"

// ErrorStatusHTTP recursively checks err.Code and
// returns appropriate http response code depending on the error,
// otherwise it returns 500
func ErrorStatusHTTP(err error) int {
	if err == nil {
		return 500
	} else if e, ok := err.(Error); ok && e.Code != "" {
		return CodeMapHTTP[e.Code]
	} else if ok && e.Err != nil {
		return ErrorStatusHTTP(e.Err)
	}
	return 500
}

// ErrorStatusGRPC recursively checks err.Code and
// returns appropriate grpc response code depending on the error,
// otherwise it returns 500
func ErrorStatusGRPC(err error) codes.Code {
	if err == nil {
		return 500
	} else if e, ok := err.(Error); ok && e.Code != "" {
		return CodeMapGRPC[e.Code]
	} else if ok && e.Err != nil {
		return ErrorStatusGRPC(e.Err)
	}
	return 500
}
