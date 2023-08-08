# Go - Asql

[![Go Report Card](https://goreportcard.com/badge/github.com/danielgatis/go-asql?style=flat-square)](https://goreportcard.com/report/github.com/danielgatis/go-asql)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/danielgatis/go-asql/master/LICENSE)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/danielgatis/go-asql)

A Go package that makes it easier to run SQL queries async

## Install

```bash
go get -u github.com/danielgatis/go-asql
```

And then import the package in your code:

```go
import "github.com/danielgatis/go-asql"
```

### Usage


```go
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
```

### License

Copyright (c) 2023-present [Daniel Gatis](https://github.com/danielgatis)

Licensed under [MIT License](./LICENSE)

### Buy me a coffee

Liked some of my work? Buy me a coffee (or more likely a beer)

<a href="https://www.buymeacoffee.com/danielgatis" target="_blank"><img src="https://bmc-cdn.nyc3.digitaloceanspaces.com/BMC-button-images/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;"></a>
