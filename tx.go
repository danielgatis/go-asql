package asql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

// Tx is a wrapper around sql.Tx that provides asynchronous methods.
type Tx struct {
	*sql.Tx

	wg sync.WaitGroup
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	tx.wg.Wait()

	return tx.Tx.Commit()
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	tx.wg.Wait()

	return tx.Tx.Rollback()
}

// Exec executes a query without returning any rows.
func (tx *Tx) Exec(query string, args ...any) (<-chan *Result, context.CancelFunc) {
	return tx.ExecContext(context.Background(), query, args...)
}

// ExecContext executes a query without returning any rows.
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...any) (<-chan *Result, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *Result)

	tx.wg.Add(1)

	go func() {
		result, err := tx.Tx.ExecContext(ctx, query, args...)

		ch <- &Result{
			Result: result,
			err:    fmt.Errorf("error executing the query: %w", err),
		}

		tx.wg.Done()
	}()

	return ch, cancel
}

// Prepare creates a prepared statement for later queries or executions.
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	return tx.PrepareContext(context.Background(), query)
}

// PrepareContext creates a prepared statement for later queries or executions.
func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := tx.Tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing the query: %w", err)
	}

	return &Stmt{
		Stmt: stmt,
		wg:   &tx.wg,
	}, nil
}

// Query executes a query that returns rows, typically a SELECT.
func (tx *Tx) Query(query string, args ...any) (<-chan *Rows, context.CancelFunc) {
	return tx.QueryContext(context.Background(), query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...any) (<-chan *Rows, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *Rows)

	tx.wg.Add(1)

	go func() {
		rows, err := tx.Tx.QueryContext(ctx, query, args...)

		ch <- &Rows{
			Rows: rows,
			err:  fmt.Errorf("error executing the query: %w:", err),
		}

		tx.wg.Done()
	}()

	return ch, cancel
}

// QueryRow executes a query that is expected to return at most one row.
func (tx *Tx) QueryRow(query string, args ...any) (<-chan *Row, context.CancelFunc) {
	return tx.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
func (tx *Tx) QueryRowContext(
	ctx context.Context,
	query string,
	args ...any,
) (<-chan *Row, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *Row)

	tx.wg.Add(1)

	go func() {
		row := tx.Tx.QueryRowContext(ctx, query, args...)

		ch <- &Row{
			Row: row,
		}

		tx.wg.Done()
	}()

	return ch, cancel
}
