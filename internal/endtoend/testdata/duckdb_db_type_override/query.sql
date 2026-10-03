-- name: CreateRecord :one
INSERT INTO records (id, data_json, preferences_json)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListRecords :many
SELECT * FROM records ORDER BY id;

-- name: ListDirectJSON :many
SELECT * FROM direct_json_view ORDER BY id;

-- name: ListCastJSON :many
SELECT * FROM cast_json_view ORDER BY id;

-- name: ListNestedJSON :many
SELECT * FROM nested_json_view ORDER BY id;

-- name: SelectJSONCasts :many
SELECT CAST(required_json AS JSON) AS required_json,
       CAST(optional_json AS JSON) AS optional_json,
       CAST(required_json AS VARCHAR) AS plain_text
FROM reviews;
