-- name: GetAlbums :many
-- SELECT
--     artists.name AS artist,
--     albums.title AS album

SELECT artists.name AS artist, albums.title AS album, year, rating
FROM artists
INNER JOIN albums_artists
    ON artists.id = albums_artists.artist_id
INNER JOIN albums
    ON albums_artists.album_id = albums.id
WHERE name LIKE ?
-- ORDER BY random() LIMIT 1
