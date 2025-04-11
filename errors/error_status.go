package errors

// ErrorStatus recursively checks err.Code and
// returns appropriate http response code depending on the error,
// otherwise it returns 500
func ErrorStatus(err error) int {
	if err == nil {
		return 500
	} else if e, ok := err.(Error); ok && e.Code != "" {
		return CodeMap[e.Code]
	} else if ok && e.Err != nil {
		return ErrorStatus(e.Err)
	}
	return 500
}
