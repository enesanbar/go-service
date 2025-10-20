package errors

// ErrorData recursively checks for data field in nested errors, if available.
func ErrorData(err error) interface{} {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Data != nil {
		return e.Data
	} else if ok && e.Err == nil {
		return ""
	} else if ok && e.Err != nil {
		return ErrorData(e.Err)
	}
	return ""
}
