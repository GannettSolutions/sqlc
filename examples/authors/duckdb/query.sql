-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: FindAuthorsByName :many
SELECT * FROM authors
WHERE name = sqlc.arg(name)
ORDER BY id;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING *;

-- name: CreateAuthors :batchexec
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
);

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;
