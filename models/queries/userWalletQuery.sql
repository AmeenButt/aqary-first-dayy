-- name: GetUserWallet :one
SELECT *
FROM user_wallet u
JOIN users us ON u.user_id = us.id
WHERE u.user_id = $1
LIMIT 1;

-- name: GetUserWalletByID :one
SELECT * FROM user_wallet u JOIN users us ON u.user_id = us.id WHERE u.id = $1 LIMIT 1;

-- name: ListUserWallets :many
SELECT * FROM user_wallet u JOIN users us ON u.user_id = us.id ORDER BY u.id;

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




