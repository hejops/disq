package main

import (
	"bufio"
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
	"golang.org/x/sys/unix"

	"disq/sqlc"
)

// TODO: consider https://github.com/cvilsmeier/sqinn-go (go-sqlite3 is cgo,
// and thus very slow to compile)

//go:embed schema.sql
var schema string

var artist = flag.String("artist", "", "")

func artistMenu(q *sqlc.Queries, ctx context.Context, s string) string {
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
	// fmt stops at first whitespace!
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		return line
	}
	return ""
}

func main() {
	// https://docs.sqlc.dev/en/stable/tutorials/getting-started-sqlite.html

	// TODO: currently dump is performed by ../gripts/disq (and written
	// there), so we rely on this file being a hardlink. in future, the
	// dumping code should be done here
	ex, _ := os.Executable()
	abs := filepath.Join(filepath.Dir(ex), "collection2.db")

	db, err := sql.Open("sqlite3", abs) // not sqlite!
	if err != nil {
		panic(err)
	}

	// create tables
	ctx := context.Background()
	if _, err = db.ExecContext(ctx, schema); err != nil {
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
	var opts []tea.ProgramOption

	// https://github.com/mattn/go-isatty/blob/master/isatty_tcgets.go
	_, err = unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETS)
	isTerm := err == nil
	if !isTerm {
		opts = append(opts, tea.WithoutRenderer())
	}
	m.interactive = isTerm

	if _, err := tea.NewProgram(&m, opts...).Run(); err != nil {
		panic(err)
	}

	// default:
	// 	fmt.Println("noop")
	// }
}
