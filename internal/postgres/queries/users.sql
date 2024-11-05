-- name: CreateUser :one
INSERT INTO users(full_name, username, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListUsers :many
SELECT id, full_name, username, is_admin FROM users ORDER BY created_at DESC;

-- name: ListUsersCount :one
SELECT count(1) FROM users;

-- name: ListUsersPublic :many
SELECT full_name, username, about_me FROM users WHERE is_banned = false;

-- name: ListUsersPublicCount :one
SELECT count(1) FROM users WHERE is_banned = false;

-- name: GetUser :one
SELECT full_name, username, about_me, created_at FROM users WHERE id = $1;

-- name: GetUserFull :one
SELECT full_name, username, about_me, is_admin, is_banned, created_at FROM users WHERE id = $1;

-- name: GetUserPublic :one
SELECT full_name, about_me FROM users WHERE username = $1;

-- name: GetUserAuthData :one
SELECT full_name, password, is_banned FROM users WHERE username = $1;

-- name: GetUserPassword :one
SELECT password FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET 
    full_name = $2,
    username = $3,
    about_me = $4
WHERE id = $1;

-- name: UpdateUserFull :exec
UPDATE users
SET 
    full_name = $2,
    username = $3,
    about_me = $4,
    password = $5,
    is_banned = $6
WHERE id = $1;