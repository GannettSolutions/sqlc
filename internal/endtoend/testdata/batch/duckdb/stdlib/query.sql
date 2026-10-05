-- name: InsertID :batchexec
INSERT INTO authors (id, name) VALUES ($1, 'unknown');

-- name: InsertAuthor :batchexec
INSERT INTO authors (id, name) VALUES ($1, $2);
