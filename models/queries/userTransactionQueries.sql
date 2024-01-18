

-- name: GetUserWalletTransactions :many
SELECT * FROM user_transactions ut JOIN user_wallet us ON ut.user_wallet_id = us.id JOIN users u ON us.user_id = u.id WHERE ut.user_wallet_id = $1;

-- name: ListTransactions :many
SELECT * FROM user_transactions u JOIN user_wallet us ON u.user_wallet_id = us.id ORDER BY u.id;

-- name: CreateUserTransaction :one
INSERT INTO user_transactions (
  user_wallet_id, transaction_amount, action
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteUserTransaction :exec
DELETE FROM user_transactions
WHERE id = $1;

