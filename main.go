package main

import (
	"context"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
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

	// TODO: fallthrough cases with squirrel?
	switch {
	case *artist != "":
		albums, err := q.GetAlbums(ctx, *artist)
		if err != nil {
			panic(err)
		}

		if albums == nil {
			return
		}

		lf, _ := tea.LogToFile("/tmp/disq.log", "")
		defer lf.Close()

		m := Model{albums: albums}
		_, err = tea.NewProgram(&m).Run()
		if err != nil {
			panic(err)
		}

	default:
		fmt.Println("noop")
	}
}
