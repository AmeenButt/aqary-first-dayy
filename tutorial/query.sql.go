// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package tutorial

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING id, name, bio
`

type CreateAuthorParams struct {
	Name string
	Bio  pgtype.Text
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	row := q.db.QueryRow(ctx, createAuthor, arg.Name, arg.Bio)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  name, email, password
) VALUES (
  $1, $2, $3
)
RETURNING id, name, email, password, author_status, author_id
`

type CreateUserParams struct {
	Name     string
	Email    pgtype.Text
	Password pgtype.Text
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Name, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.AuthorStatus,
		&i.AuthorID,
	)
	return i, err
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteAuthor, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, bio FROM authors
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAuthor(ctx context.Context, id int64) (Author, error) {
	row := q.db.QueryRow(ctx, getAuthor, id)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const getUsers = `-- name: GetUsers :one
SELECT id, name, email, password, author_status, author_id FROM users AS u WHERE u.id = $1 LIMIT 1
`

func (q *Queries) GetUsers(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUsers, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.AuthorStatus,
		&i.AuthorID,
	)
	return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, name, bio FROM authors
ORDER BY name
`

func (q *Queries) ListAuthors(ctx context.Context) ([]Author, error) {
	rows, err := q.db.Query(ctx, listAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Name, &i.Bio); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT id, name, email, password, author_status, author_id FROM users ORDER BY name
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.AuthorStatus,
			&i.AuthorID,
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

const updateAuthor = `-- name: UpdateAuthor :exec
UPDATE authors
  set name = $2,
  bio = $3
WHERE id = $1 RETURNING id, name, bio
`

type UpdateAuthorParams struct {
	ID   int64
	Name string
	Bio  pgtype.Text
}

func (q *Queries) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) error {
	_, err := q.db.Exec(ctx, updateAuthor, arg.ID, arg.Name, arg.Bio)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
  set name = $2,
  email = $3,
  password = $4
WHERE id = $1 RETURNING id, name, email, password, author_status, author_id
`

type UpdateUserParams struct {
	ID       int64
	Name     string
	Email    pgtype.Text
	Password pgtype.Text
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Password,
	)
	return err
}