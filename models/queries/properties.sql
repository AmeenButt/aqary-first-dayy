-- name: InsertProperty :one
INSERT INTO properties (sizeInSqFeet, location,demand,status,user_id) VALUES ($1,$2,$3,$4, $5) RETURNING *;

-- name: UpdateProperty :one
UPDATE properties SET sizeInSqFeet =$2 , location = $3 ,demand = $4 ,status=$5 
WHERE id=$1 RETURNING *;

-- name: GetPropertyByID :one
SELECT p.*, us.* FROM properties p JOIN users us ON p.user_id = us.id WHERE p.id = $1;

-- name: GetPropertyByUserID :many
SELECT p.*, us.* FROM properties p JOIN users us ON p.user_id = us.id WHERE p.user_id = $1;

-- name: UpdateStatus :one
UPDATE properties SET status =$2  
WHERE id=$1 RETURNING *;

-- name: DeleteProperty :one
DELETE FROM properties WHERE id=$1 RETURNING *;