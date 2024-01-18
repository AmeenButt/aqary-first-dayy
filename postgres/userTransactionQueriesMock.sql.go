package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (q *Queries) TestCreateUserTransaction(ctx context.Context, arg CreateUserTransactionParams) (UserTransaction, error) {
	var i UserTransaction
	return i, nil
}

func (q *Queries) TestDeleteUserTransaction(ctx context.Context, id int64) error {
	_, nil := q.db.Exec(ctx, deleteUserTransaction, id)
	return nil
}

func (q *Queries) TestGetUserWalletTransactions(ctx context.Context, userWalletID pgtype.Int4) ([]GetUserWalletTransactionsRow, error) {
	var items []GetUserWalletTransactionsRow
	return items, nil
}

func (q *Queries) TestListTransactions(ctx context.Context) ([]ListTransactionsRow, error) {
	var items []ListTransactionsRow
	return items, nil
}
