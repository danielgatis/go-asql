package asql

import "database/sql"

// Result is a wrapper around sql.Result that provides an Err method.
type Result struct {
	sql.Result

	err error
}

// Err returns the error, if any, that was encountered during the query.
func (r *Result) Err() error {
	return r.err
}
