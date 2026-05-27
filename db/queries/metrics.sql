-- name: FindById :one
SELECT id, ts, type, name, delta, value
FROM metrics
WHERE id = $1
LIMIT 1;

-- name: FindByName :many
SELECT id, ts, type, name, delta, value
FROM metrics
WHERE name = $1;

-- name: FindAll :many
SELECT id, ts, type, name, delta, value
FROM metrics;

-- name: Insert :one
INSERT INTO metrics(ts, type, name, delta, value)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, ts, type, name, delta, value;

-- name: InsertBatch :copyfrom
INSERT INTO metrics(ts, type, name, delta, value)
VALUES ($1, $2, $3, $4, $5);

-- name: Update :one
UPDATE metrics
SET ts    = $2,
    type  = $3,
    name  = $4,
    delta = $5,
    value = $6
WHERE id = $1
RETURNING id, ts, type, name, delta, value;

-- name: DeleteById :execrows
DELETE
FROM metrics
WHERE id = $1;
