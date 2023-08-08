package asql_test

import (
	"context"
	"testing"

	"github.com/danielgatis/go-asql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestTxCommit(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	tx, _ := db.BeginTx(context.Background(), nil)
	rc1, _ := tx.Exec(`insert into test_table (id, name) values (3, "jack")`)
	rc2, _ := tx.Exec(`insert into test_table (id, name) values (4, "jill")`)

	<-rc1
	<-rc2

	tx.Commit()

	rc, _ := db.QueryRow(`select count(*) from test_table`)
	row := <-rc

	var count int
	row.Scan(&count)

	assert.Equal(t, 4, count)
}

func TestTxRollback(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	tx, _ := db.BeginTx(context.Background(), nil)
	rc1, _ := tx.Exec(`insert into test_table (id, name) values (3, "jack")`)
	rc2, _ := tx.Exec(`insert into test_table (id, name) values (4, "jill")`)

	<-rc1
	<-rc2

	tx.Rollback()

	rc, _ := db.QueryRow(`select count(*) from test_table`)
	row := <-rc

	var count int
	row.Scan(&count)

	assert.Equal(t, 2, count)
}

func TestTxExec(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	tx, _ := db.BeginTx(context.Background(), nil)
	query := `insert into test_table (id, name) values (3, "jack")`
	rc, _ := tx.Exec(query)
	result := <-rc

	actual, _ := result.RowsAffected()

	assert.EqualValues(t, 1, actual)
}

func TestTxPrepare(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	tx, _ := db.BeginTx(context.Background(), nil)
	stmt, _ := tx.Prepare(`insert into test_table (id, name) values (?, ?)`)

	rc, _ := stmt.Exec(4, "jill")
	result := <-rc

	actual, _ := result.RowsAffected()
	assert.EqualValues(t, 1, actual)
}

func TestTxQuery(t *testing.T) {
	t.Parallel()

	type TestTable struct {
		ID   int
		Name string
	}

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	tx, _ := db.BeginTx(context.Background(), nil)

	rc, _ := tx.Query(`select * from test_table`)
	rows := <-rc

	records := make([]TestTable, 0)

	for rows.Next() {
		var record TestTable
		rows.Scan(&record.ID, &record.Name)
		records = append(records, record)
	}

	expected := []TestTable{
		{1, "alice"},
		{2, "bob"},
	}

	assert.EqualValues(t, expected, records)
}

func TestTxQueryRow(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	tx, _ := db.BeginTx(context.Background(), nil)
	rc, _ := tx.QueryRow(`select count(*) from test_table`)
	row := <-rc

	var count int
	row.Scan(&count)

	assert.Equal(t, count, 2)
}
