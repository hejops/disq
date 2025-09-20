package main

import (
	"fmt"
	"iter"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"disq/sqlc"
)

type Model struct {
	t           table.Model
	albums      []sqlc.GetAlbumsRow
	interactive bool
	done        bool
}

func must[T any](x T, _ error) T { return x }

func asRow(album sqlc.GetAlbumsRow) table.Row {
	y := must(album.Year.Value()).(int64)
	r := must(album.Rating.Value()).(int64)
	return table.Row{
		album.Artist,
		album.Album,
		strconv.Itoa(int(y)),
		strconv.Itoa(int(r)),
	}
}

func maxLen(it iter.Seq[string]) int {
	var w int
	for s := range it {
		w = max(w, len(s))
	}
	return w
}

func (m *Model) getAlbums() iter.Seq[string] { // meh
	return func(yield func(string) bool) {
		for _, album := range m.albums {
			if !yield(album.Album) {
				return
			}
		}
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *Model) Init() tea.Cmd {
	var rows []table.Row
	for _, album := range m.albums {
		rows = append(rows, asRow(album))
	}
	log.Println(len(rows), "rows")

	m.t = table.New() // this instantiation makes struct embedding tricky
	m.t.SetColumns([]table.Column{
		{Title: "Artist", Width: len(m.albums[0].Artist)},
		{Title: "Album", Width: maxLen(m.getAlbums())}, // TODO: truncate
		{Title: "Year", Width: 4},
		{Title: "Rating", Width: 6},
	})
	m.t.SetRows(rows)

	// if large number of rows, need pager behaviour, so don't exit
	if len(m.t.Rows()) > 100 {
		return nil
	} else {
		return tea.Quit
	}
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch s {
		case "x":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *Model) View() string {
	s := lipgloss.NewStyle().
		MaxHeight(len(m.albums) + 2).
		Render(m.t.View())
	if !m.interactive && !m.done {
		// [ -t 1 ] false -> WithoutRenderer -> no Render -> need to print ourselves
		fmt.Println(s)
		m.done = true // View is called 2x on program init, for some reason
	}
	return s
}
