package main

import (
	"context"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"

	"disq/sqlc"
)

//go:embed schema.sql
var schema string

var artist = flag.String("artist", "", "")

func artistMenu(
	q *sqlc.Queries,
	ctx context.Context,
	s string,
) string {
	artists, err := q.GetArtistsWithSubstring(ctx, sql.NullString{String: s, Valid: true})
	if err != nil {
		panic(err)
	}

	for i, a := range artists {
		fmt.Println(i+1, a)
	}

	switch len(artists) {
	case 0:
		// panic("no artist")
		os.Exit(1)
		return ""
	case 1:
		return artists[0]
	default:
		n := readLine()
		if n == "" {
			// artistMenu(q, ctx, s)
			return artists[0]
		}
		return artists[must(strconv.Atoi(n))-1]
	}
}

func readLine() string {
	var dest string
	_, _ = fmt.Scanln(&dest)
	// if err != nil { // unexpected newline
	// 	panic(err)
	// }
	return dest
}

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

	if *artist == "" {
		fmt.Print("search artist: ")
		search := readLine()
		if search == "" {
			return
		}
		a := artistMenu(q, ctx, search)
		artist = &a
	}

	// TODO: fallthrough cases with squirrel?
	// switch {
	// case *artist != "":
	albums, err := q.GetAlbums(ctx, *artist)
	if err != nil {
		panic(err)
	}

	if albums == nil {
		fmt.Println("no results")
		return
	}

	lf, _ := tea.LogToFile("/tmp/disq.log", "")
	defer lf.Close()

	fmt.Println(len(albums), "albums")
	m := Model{albums: albums}
	_, err = tea.NewProgram(&m).Run()
	if err != nil {
		panic(err)
	}

	// default:
	// 	fmt.Println("noop")
	// }
}
