package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (q *Queries) TestCreateUserWallet(ctx context.Context, arg CreateUserWalletParams) (UserWallet, error) {
	var i UserWallet
	return i, nil
}

func (q *Queries) TestDeleteUserWallet(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUserWallet, id)
	return err
}
func (q *Queries) TestGetUserWallet(ctx context.Context, userID pgtype.Int4) (GetUserWalletRow, error) {
	var i GetUserWalletRow
	return i, nil
}

func (q *Queries) TestGetUserWalletByID(ctx context.Context, id int64) (GetUserWalletByIDRow, error) {
	var i GetUserWalletByIDRow
	return i, nil
}

func (q *Queries) TestListUserWallets(ctx context.Context) ([]ListUserWalletsRow, error) {
	var items []ListUserWalletsRow
	return items, nil
}
func (q *Queries) TestUpdateUserWalletAmount(ctx context.Context, arg UpdateUserWalletAmountParams) error {
	_, err := q.db.Exec(ctx, updateUserWalletAmount, arg.ID, arg.Amount)
	return err
}
