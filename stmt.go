package asql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

// Stmt is a wrapper around sql.Stmt that provides asynchronous methods.
type Stmt struct {
	*sql.Stmt

	wg *sync.WaitGroup
}

// Exec executes a prepared statement with the given arguments without returning any rows.
func (s *Stmt) Exec(args ...any) (<-chan *Result, context.CancelFunc) {
	return s.ExecContext(context.Background(), args...)
}

// ExecContext executes a prepared statement with the given arguments without returning any rows.
func (s *Stmt) ExecContext(ctx context.Context, args ...any) (<-chan *Result, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *Result)

	if s.wg != nil {
		s.wg.Add(1)
	}

	go func() {
		result, err := s.Stmt.ExecContext(ctx, args...)

		ch <- &Result{
			Result: result,
			err:    fmt.Errorf("error executing the query: %w", err),
		}

		if s.wg != nil {
			s.wg.Done()
		}
	}()

	return ch, cancel
}

// Query executes a prepared query statement with the given arguments.
func (s *Stmt) Query(args ...any) (<-chan *Rows, context.CancelFunc) {
	return s.QueryContext(context.Background(), args...)
}

// QueryContext executes a prepared query statement with the given arguments.
func (s *Stmt) QueryContext(ctx context.Context, args ...any) (<-chan *Rows, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *Rows)

	if s.wg != nil {
		s.wg.Add(1)
	}

	go func() {
		rows, err := s.Stmt.QueryContext(ctx, args...)

		ch <- &Rows{
			Rows: rows,
			err:  fmt.Errorf("error querying the database: %w", err),
		}

		if s.wg != nil {
			s.wg.Done()
		}
	}()

	return ch, cancel
}

// QueryRow executes a prepared query statement with the given arguments.
func (s *Stmt) QueryRow(args ...any) (<-chan *Row, context.CancelFunc) {
	return s.QueryRowContext(context.Background(), args...)
}

// QueryRowContext executes a prepared query statement with the given arguments.
func (s *Stmt) QueryRowContext(ctx context.Context, args ...any) (<-chan *Row, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *Row)

	if s.wg != nil {
		s.wg.Add(1)
	}

	go func() {
		row := s.Stmt.QueryRowContext(ctx, args...)

		ch <- &Row{
			Row: row,
		}

		if s.wg != nil {
			s.wg.Done()
		}
	}()

	return ch, cancel
}
