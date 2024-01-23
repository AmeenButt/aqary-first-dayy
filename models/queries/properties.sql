-- name: InsertProperty :one
INSERT INTO properties (sizeInSqFeet, location,demand,user_id,images) VALUES ($1,$2,$3,$4,$5) RETURNING *;

-- name: UpdateProperty :one
UPDATE properties SET sizeInSqFeet =$2 , location = $3 ,demand = $4 ,status=$5 , images=$6,
  updated_at = NOW()
WHERE id=$1 RETURNING *;

-- name: GetPropertyByID :one
SELECT p.*, us.* FROM properties p JOIN users us ON p.user_id = us.id WHERE p.id = $1;

-- name: GetPropertyByUserID :many
SELECT p.*, us.* FROM properties p JOIN users us ON p.user_id = us.id WHERE p.user_id = $1;

-- name: UpdateStatus :one
UPDATE properties SET status =$2  ,
  updated_at = NOW()
WHERE id=$1 RETURNING *;

-- name: DeleteProperty :one
DELETE FROM properties WHERE id=$1 RETURNING *;


-- name: UpdatePropertySize :one
UPDATE properties SET sizeInSqFeet =$2,
  updated_at = NOW()
WHERE id=$1 RETURNING *;

-- name: UpdatePropertyLocation :one
UPDATE properties SET location =$2,
  updated_at = NOW()
WHERE id=$1 RETURNING *;

-- name: UpdatePropertyDemand :one
UPDATE properties SET demand =$2,
  updated_at = NOW()
WHERE id=$1 RETURNING *;

-- name: UpdatePropertyImages :one
UPDATE properties SET images =$2,
  updated_at = NOW()
WHERE id=$1 RETURNING *;