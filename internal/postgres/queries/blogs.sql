-- name: CreateBlog :one
INSERT INTO blogs(author_id, title, summary, content)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListBlogs :many
SELECT
    b.id,
    b.title,
    b.removed_at is not null as is_removed,
    u.username as author_username
FROM blogs b
JOIN users u ON u.id = b.author_id
ORDER BY b.created_at DESC;

-- name: ListBlogsCount :one
SELECT count(1) FROM blogs;

-- name: ListBlogsPublic :many
SELECT
    b.id,
    b.title,
    b.summary,
    u.full_name as author_name,
    b.created_at as published_at
FROM blogs b
JOIN users u ON u.id = b.author_id
WHERE 
    b.removed_at is null
    and u.is_banned = false
ORDER BY b.created_at;

-- name: ListBlogsPublicCount :one
SELECT count(1) 
FROM blogs b
JOIN users u ON u.id = b.author_id
WHERE 
    b.removed_at is null
    and u.is_banned = false;

-- name: ListAuthorBlogs :many
SELECT
    b.id,
    b.title,
    b.summary,
    b.created_at as published_at
FROM blogs b
JOIN users u ON u.id = b.author_id
WHERE 
    b.removed_at is null
    and u.is_banned = false
    and u.username = $1
ORDER BY b.created_at;

-- name: GetBlog :one
SELECT * FROM blogs WHERE id = $1;

-- name: GetBlogPublic :one
SELECT
    b.title,
    b.summary,
    b.content,
    b.created_at,
    b.updated_at,
    u.full_name as author_name,
    u.username  as author_username,
    u.about_me  as author_about
FROM blogs b
JOIN users u ON u.id = b.author_id
WHERE
    b.id = $1 and 
    b.removed_at is null and
    u.is_banned = false;

-- name: UpdateBlog :exec
UPDATE blogs
SET 
    title = $2,
    summary = $3,
    content = $4,
    updated_at = NOW()
WHERE id = $1;

-- name: RemoveBlog :exec
UPDATE blogs
SET removed_at = NOW()
WHERE id = $1;