-- name: GetUserWallet :one
SELECT * FROM user_wallet AS u WHERE u.id = $1 LIMIT 1;

-- name: ListUserWallets :many
SELECT * FROM user_wallet ORDER BY id;

-- name: CreateUserWallet :one
INSERT INTO user_wallet (
  user_id, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateUserWalletAmount :exec
UPDATE user_wallet
  set amount = $2
  WHERE id = $1 RETURNING *;

-- name: DeleteUserWallet :exec
DELETE FROM user_wallet
WHERE id = $1;

-- name: GetUserWalletTransactions :one
SELECT * FROM user_transactions WHERE user_wallet_id = $1;

-- name: ListTransactions :many
SELECT * FROM user_transactions ORDER BY id;

-- name: CreateUserTransaction :one
INSERT INTO user_transactions (
  user_wallet_id, transaction_amount
) VALUES (
  $1, $2
)
RETURNING *;


-- name: DeleteUserTransaction :exec
DELETE FROM user_transactions
WHERE id = $1;

-- name: GetUsers :one
SELECT * FROM users AS u WHERE u.id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  name, email, password
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
  set name = $2,
  email = $3,
  password = $4
WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;