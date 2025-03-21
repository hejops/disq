// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: artists.sql

package sqlc

import (
	"context"
	"database/sql"
)

const getAlbums = `-- name: GetAlbums :many
SELECT artists.name AS artist, albums.title AS album, year, rating
FROM artists
INNER JOIN albums_artists
    ON artists.id = albums_artists.artist_id
INNER JOIN albums
    ON albums_artists.album_id = albums.id
WHERE name LIKE ?
ORDER BY rating DESC
`

type GetAlbumsRow struct {
	Artist string
	Album  string
	Year   sql.NullInt64
	Rating sql.NullInt64
}

func (q *Queries) GetAlbums(ctx context.Context, name string) ([]GetAlbumsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAlbums, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAlbumsRow
	for rows.Next() {
		var i GetAlbumsRow
		if err := rows.Scan(
			&i.Artist,
			&i.Album,
			&i.Year,
			&i.Rating,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getArtistsWithSubstring = `-- name: GetArtistsWithSubstring :many
SELECT name --AS artist
FROM artists
WHERE name LIKE '%' || ? || '%'
`

func (q *Queries) GetArtistsWithSubstring(ctx context.Context, dollar_1 sql.NullString) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getArtistsWithSubstring, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
