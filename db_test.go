package asql_test

import (
	"testing"

	"github.com/danielgatis/go-asql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestDBOpen(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")

	assert.NotNil(t, db)
}

func TestDBClose(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	assert.Nil(t, db.Close())
}

func TestDBPing(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	assert.Nil(t, db.Ping())
}

func TestDBBegin(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	tx, _ := db.Begin()

	assert.NotNil(t, tx)
}

func TestDBExec(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	query := `insert into test_table (id, name) values (3, "jack")`
	rc, _ := db.Exec(query)
	result := <-rc

	actual, err := result.RowsAffected()
	if err != nil {
		t.Error(err)
		return
	}

	assert.EqualValues(t, 1, actual)
}

func TestDBPrepare(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	stmt, _ := db.Prepare(`insert into test_table (id, name) values (?, ?)`)
	rc, _ := stmt.Exec(4, "jill")

	result := <-rc
	actual, _ := result.RowsAffected()

	assert.EqualValues(t, 1, actual)
}

func TestDBQuery(t *testing.T) {
	t.Parallel()

	type TestTable struct {
		ID   int
		Name string
	}

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	rc, _ := db.Query(`select * from test_table`)
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

func TestDBQueryRow(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	rc, _ := db.QueryRow(`select count(*) from test_table`)
	row := <-rc

	var count int
	row.Scan(&count)

	assert.Equal(t, count, 2)
}
