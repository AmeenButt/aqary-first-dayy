package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  name, email, password, profile_picture
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, name, email, password, profile_picture, created_at, updated_at
`

type CreateUserParams struct {
	Name           string      `json:"name"`
	Email          pgtype.Text `json:"email"`
	Password       pgtype.Text `json:"password"`
	ProfilePicture pgtype.Text `json:"profile_picture"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.ProfilePicture,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.ProfilePicture,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, password, profile_picture, created_at, updated_at FROM users WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.ProfilePicture,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, email, password, profile_picture, created_at, updated_at FROM users WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.ProfilePicture,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, name, email, password, profile_picture, created_at, updated_at FROM users ORDER BY name
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.ProfilePicture,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
  set name = $2,
  email = $3,
  password = $4,
  profile_picture = $5,
  updated_at = NOW()
WHERE id = $1 RETURNING id, name, email, password, profile_picture, created_at, updated_at
`

type UpdateUserParams struct {
	ID             int64       `json:"id"`
	Name           string      `json:"name"`
	Email          pgtype.Text `json:"email"`
	Password       pgtype.Text `json:"password"`
	ProfilePicture pgtype.Text `json:"profile_picture"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.ProfilePicture,
	)
	return err
}

const updateUserPicture = `-- name: UpdateUserPicture :exec
UPDATE users
  set profile_picture = $2,
  updated_at = NOW()
WHERE id = $1 RETURNING id, name, email, password, profile_picture, created_at, updated_at
`

type UpdateUserPictureParams struct {
	ID             int64       `json:"id"`
	ProfilePicture pgtype.Text `json:"profile_picture"`
}

func (q *Queries) UpdateUserPicture(ctx context.Context, arg UpdateUserPictureParams) error {
	_, err := q.db.Exec(ctx, updateUserPicture, arg.ID, arg.ProfilePicture)
	return err
}
