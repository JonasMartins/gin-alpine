-- name: SearchUsers :many
WITH users_page AS (
   SELECT u.id,
      u.uuid,
      u.name,
      u.email,
      r.role,
      u.enabled,
      u.created_at,
      u.updated_at
   FROM users u
      LEFT JOIN roles r ON u.role_id = r.id
   WHERE u.deleted_at IS NULL
      AND (
         LOWER(u.name) LIKE '%' || $1 || '%'
         OR LOWER(u.email) LIKE '%' || $2 || '%'
      )
   ORDER BY u.id DESC
   LIMIT $3 OFFSET $4
),
total AS (
   SELECT COUNT(*) AS total_count
   FROM users u
   WHERE u.deleted_at IS NULL
      AND (
         LOWER(u.name) LIKE '%' || $5 || '%'
         OR LOWER(u.email) LIKE '%' || $6 || '%'
      )
)
SELECT up.id,
   up.uuid,
   up.name,
   up.email,
   up.role,
   up.enabled,
   up.created_at,
   up.updated_at,
   t.total_count
FROM users_page up
   CROSS JOIN total t;
-- name: CreateUser :one
INSERT INTO users (uuid, name, email, password, role_id, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
-- name: GetUsersWithTotal :many
WITH users_page AS (
   SELECT u.id,
      u.uuid,
      u.name,
      u.email,
      r.role,
      u.enabled,
      u.created_at,
      u.updated_at
   FROM users u
      LEFT JOIN roles r ON u.role_id = r.id
   WHERE u.deleted_at IS NULL
   ORDER BY u.id DESC
   LIMIT $1 OFFSET $2
),
total AS (
   SELECT COUNT(*) AS total_count
   FROM users u
   WHERE u.deleted_at IS NULL
)
SELECT up.id,
   up.uuid,
   up.name,
   up.email,
   up.role,
   up.enabled,
   up.created_at,
   up.updated_at,
   t.total_count
FROM users_page up
   CROSS JOIN total t;
-- name: GetTotalUsers :one
SELECT COUNT(*)
FROM users u
WHERE u.deleted_at IS NULL;
-- name: UpdateUser :exec
UPDATE users
SET name = COALESCE(sqlc.narg('name'), name),
   email = COALESCE(sqlc.narg('email'), email),
   password = COALESCE(sqlc.narg('password'), password),
   updated_at = sqlc.arg('updated_at')
WHERE id = sqlc.arg('id');
-- name: UpdateUserAdmin :exec
UPDATE users
SET name = COALESCE(sqlc.narg('name'), name),
   email = COALESCE(sqlc.narg('email'), email),
   password = COALESCE(sqlc.narg('password'), password),
   updated_at = sqlc.arg('updated_at'),
   role_id = COALESCE(sqlc.narg('role_id'), role_id),
   deleted_at = COALESCE(sqlc.narg('deleted_at'), deleted_at),
   enabled = COALESCE(sqlc.narg('enabled'), enabled)
WHERE id = sqlc.arg('id');
-- name: FindUserByEmail :one
SELECT u.id,
   u.uuid,
   u.name,
   u.email,
   r.role,
   u.role_id,
   u.enabled,
   u.password,
   u.created_at,
   u.updated_at
FROM users u
   LEFT JOIN roles r ON u.role_id = r.id
WHERE u.email = $1
   AND u.deleted_at IS NULL
LIMIT 1;
-- name: FindUserByID :one
SELECT u.id,
   u.uuid,
   u.name,
   u.email,
   r.role,
   u.enabled,
   u.password,
   u.created_at,
   u.updated_at
FROM users u
   LEFT JOIN roles r ON u.role_id = r.id
WHERE u.id = $1
   AND u.deleted_at IS NULL
LIMIT 1;
-- name: FindRoleByID :one
SELECT *
FROM roles r
WHERE r.id = $1;