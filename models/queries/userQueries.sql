
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  name, email, password, profile_picture
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
  set name = $2,
  email = $3,
  password = $4,
  profile_picture = $5,
  updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateOTP :one
UPDATE users
  set otp = $2,
  updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserPicture :exec
UPDATE users
  set profile_picture = $2,
  updated_at = NOW()
WHERE id = $1 RETURNING *;