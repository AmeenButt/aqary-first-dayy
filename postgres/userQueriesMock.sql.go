package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type TestCreateUserParams struct {
	Name           string
	Email          pgtype.Text
	Password       pgtype.Text
	ProfilePicture pgtype.Text
}

func (q *Queries) TestCreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	var i User
	return i, nil
}

func (q *Queries) TestGetUserByEmail(ctx context.Context, email pgtype.Text) (User, error) {
	var i User
	return i, nil
}
func (q *Queries) TestGetUserByID(ctx context.Context, id int64) (User, error) {
	var i User
	return i, nil
}
func (q *Queries) TestListUsers(ctx context.Context) ([]User, error) {
	var items []User
	return items, nil
}
func (q *Queries) TestUpdateUser(ctx context.Context, arg UpdateUserParams) error {
	return nil
}
func (q *Queries) TestUpdateUserPicture(ctx context.Context, arg UpdateUserPictureParams) error {
	return nil
}
