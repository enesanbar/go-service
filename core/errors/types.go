package errors

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
