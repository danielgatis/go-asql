package asql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

// DB is a wrapper around sql.DB that provides asynchronous methods.
type DB struct {
	*sql.DB
}

// Open opens a database specified by its database driver name and dsn.
func Open(driverName, dsn string) (*DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening the database connection: %w", err)
	}

	return &DB{
		DB: db,
	}, nil
}

// Close closes the database.
func (db *DB) Close() error {
	return db.DB.Close()
}

// Ping verifies a connection to the database is still alive.
func (db *DB) Ping() error {
	return db.DB.Ping()
}

// Load loads a query from a file and executes it.
func (db *DB) Load(path string) error {
	query, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading the file: %w", err)
	}

	_, err = db.DB.Exec(string(query))
	if err != nil {
		return fmt.Errorf("error executing the query: %w", err)
	}

	return nil
}

// BeginTx starts a transaction.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("error starting a new transaction: %w", err)
	}

	return &Tx{
		Tx: tx,
	}, nil
}

// Exec executes a query without returning any rows.
func (db *DB) Exec(query string, args ...any) (<-chan *Result, context.CancelFunc) {
	return db.ExecContext(context.Background(), query, args...)
}

// ExecContext executes a query without returning any rows.
func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (<-chan *Result, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *Result)

	go func() {
		result, err := db.DB.ExecContext(ctx, query, args...)

		ch <- &Result{
			Result: result,
			err:    fmt.Errorf("error executing the query: %w", err),
		}
	}()

	return ch, cancel
}

// Prepare creates a prepared statement for later queries or executions.
func (db *DB) Prepare(query string) (*Stmt, error) {
	return db.PrepareContext(context.Background(), query)
}

// PrepareContext creates a prepared statement for later queries or executions.
func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing the query: %w", err)
	}

	return &Stmt{
		Stmt: stmt,
	}, nil
}

// Query executes a query that returns rows, typically a SELECT.
func (db *DB) Query(query string, args ...any) (<-chan *Rows, context.CancelFunc) {
	return db.QueryContext(context.Background(), query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (<-chan *Rows, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *Rows)

	go func() {
		rows, err := db.DB.QueryContext(ctx, query, args...)

		ch <- &Rows{
			Rows: rows,
			err:  fmt.Errorf("error executing the query: %w", err),
		}
	}()

	return ch, cancel
}

// QueryRow executes a query that is expected to return at most one row.
func (db *DB) QueryRow(query string, args ...any) (<-chan *Row, context.CancelFunc) {
	return db.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
func (db *DB) QueryRowContext(
	ctx context.Context,
	query string,
	args ...any,
) (<-chan *Row, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *Row)

	go func() {
		row := db.DB.QueryRowContext(ctx, query, args...)

		ch <- &Row{
			Row: row,
		}
	}()

	return ch, cancel
}
