package main

import (
	"context"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"disq/sqlc"
)

//go:embed schema.sql
var schema string

var artist = flag.String("artist", "", "")

func main() {
	// https://docs.sqlc.dev/en/stable/tutorials/getting-started-sqlite.html

	ex, _ := os.Executable()
	abs := filepath.Join(filepath.Dir(ex), "collection2.db")

	db, err := sql.Open(
		"sqlite3", // not sqlite!
		abs,
	)
	if err != nil {
		panic(err)
	}

	// create tables
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, schema); err != nil {
		panic(err)
	}

	q := sqlc.New(db)

	flag.Parse()

	switch {
	case *artist != "":
		rows, err := q.GetAlbums(ctx, *artist)
		if err != nil {
			panic(err)
		}

		// TODO: tabular
		fmt.Println(rows)

	default:
		fmt.Println("noop")
	}
}
