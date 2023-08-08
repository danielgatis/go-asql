package asql

import "database/sql"

// Row is a wrapper around sql.Row.
type Row struct {
	*sql.Row
}

// Rows is a wrapper around sql.Rows that provides an Err method.
type Rows struct {
	*sql.Rows

	err error
}

// Err returns the error, if any, that was encountered during the query.
func (rs *Rows) Err() error {
	if rs.err != nil {
		return rs.err
	}

	return rs.Rows.Err()
}
