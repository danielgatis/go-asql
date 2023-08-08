package main

import (
	"fmt"

	"github.com/danielgatis/go-asql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
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

	fmt.Print(records)
}
