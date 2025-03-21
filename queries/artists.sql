-- name: GetAlbums :many
SELECT artists.name AS artist, albums.title AS album, year, rating
FROM artists
INNER JOIN albums_artists
    ON artists.id = albums_artists.artist_id
INNER JOIN albums
    ON albums_artists.album_id = albums.id
WHERE name LIKE ?
ORDER BY rating DESC;

-- name: GetArtist :many
SELECT name --AS artist
FROM artists
WHERE name LIKE '%' || ? || '%'; -- LIKE is case-insensitive
-- TODO: order by number of albs
