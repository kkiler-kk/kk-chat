-- name: InsertUser :execresult
INSERT INTO users (identity, name, password_hash, email) VALUES (?, ?, ?, ?);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetByIdentity :one
SELECT * FROM users WHERE identity = ?;

-- name: GetByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: UpdateUser :exec
UPDATE users SET
  name       = COALESCE(sqlc.narg(name), name),
  phone      = COALESCE(sqlc.narg(phone), phone),
  email      = COALESCE(sqlc.narg(email), email),
  avatar     = COALESCE(sqlc.narg(avatar), avatar),
  signature  = COALESCE(sqlc.narg(signature), signature),
  birth_date = COALESCE(sqlc.narg(birth_date), birth_date)
WHERE id = sqlc.arg(id);

-- name: UpdatePassword :exec
UPDATE users SET password_hash = ? WHERE id = ?;

-- name: SearchUsers :many
SELECT * FROM users
WHERE identity LIKE CONCAT(?, '%') OR name LIKE CONCAT(?, '%') OR email LIKE CONCAT(?, '%')
LIMIT ?;

-- name: GetUsersByIDs :many
SELECT * FROM users WHERE id IN (sqlc.slice(ids));
