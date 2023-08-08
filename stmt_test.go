package asql_test

import (
	"testing"

	"github.com/danielgatis/go-asql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestStmtExec(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	query := `insert into test_table (id, name) values (3, "jack")`
	stmt, _ := db.Prepare(query)
	rc, _ := stmt.Exec()
	result := <-rc

	actual, _ := result.RowsAffected()

	assert.EqualValues(t, 1, actual)
}

func TestStmtQuery(t *testing.T) {
	t.Parallel()

	type TestTable struct {
		ID   int
		Name string
	}

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	stmt, _ := db.Prepare(`select * from test_table`)
	rc, _ := stmt.Query()
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

func TestStmtQueryRow(t *testing.T) {
	t.Parallel()

	db, _ := asql.Open("sqlite3", "file::memory:")
	db.Load("testdata/schema.sql")

	stmt, _ := db.Prepare(`select count(*) from test_table`)
	rc, _ := stmt.QueryRow()
	row := <-rc

	var count int
	row.Scan(&count)

	assert.Equal(t, count, 2)
}
