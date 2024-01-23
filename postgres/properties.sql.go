// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: properties.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteProperty = `-- name: DeleteProperty :one
DELETE FROM properties WHERE id=$1 RETURNING id, sizeinsqfeet, location, images, demand, status, user_id, created_at, updated_at
`

func (q *Queries) DeleteProperty(ctx context.Context, id int64) (Property, error) {
	row := q.db.QueryRow(ctx, deleteProperty, id)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Sizeinsqfeet,
		&i.Location,
		&i.Images,
		&i.Demand,
		&i.Status,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPropertyByID = `-- name: GetPropertyByID :one
SELECT p.id, p.sizeinsqfeet, p.location, p.images, p.demand, p.status, p.user_id, p.created_at, p.updated_at, us.id, us.name, us.email, us.password, us.profile_picture, us.otp, us.created_at, us.updated_at FROM properties p JOIN users us ON p.user_id = us.id WHERE p.id = $1
`

type GetPropertyByIDRow struct {
	ID             int64            `json:"id"`
	Sizeinsqfeet   pgtype.Int4      `json:"sizeinsqfeet"`
	Location       pgtype.Text      `json:"location"`
	Images         []string         `json:"images"`
	Demand         pgtype.Text      `json:"demand"`
	Status         pgtype.Text      `json:"status"`
	UserID         pgtype.Int4      `json:"user_id"`
	CreatedAt      pgtype.Timestamp `json:"created_at"`
	UpdatedAt      pgtype.Timestamp `json:"updated_at"`
	ID_2           int64            `json:"id_2"`
	Name           string           `json:"name"`
	Email          pgtype.Text      `json:"email"`
	Password       pgtype.Text      `json:"password"`
	ProfilePicture pgtype.Text      `json:"profile_picture"`
	Otp            pgtype.Int4      `json:"otp"`
	CreatedAt_2    pgtype.Timestamp `json:"created_at_2"`
	UpdatedAt_2    pgtype.Timestamp `json:"updated_at_2"`
}

func (q *Queries) GetPropertyByID(ctx context.Context, id int64) (GetPropertyByIDRow, error) {
	row := q.db.QueryRow(ctx, getPropertyByID, id)
	var i GetPropertyByIDRow
	err := row.Scan(
		&i.ID,
		&i.Sizeinsqfeet,
		&i.Location,
		&i.Images,
		&i.Demand,
		&i.Status,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ID_2,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.ProfilePicture,
		&i.Otp,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
	)
	return i, err
}

const getPropertyByUserID = `-- name: GetPropertyByUserID :many
SELECT p.id, p.sizeinsqfeet, p.location, p.images, p.demand, p.status, p.user_id, p.created_at, p.updated_at, us.id, us.name, us.email, us.password, us.profile_picture, us.otp, us.created_at, us.updated_at FROM properties p JOIN users us ON p.user_id = us.id WHERE p.user_id = $1
`

type GetPropertyByUserIDRow struct {
	ID             int64            `json:"id"`
	Sizeinsqfeet   pgtype.Int4      `json:"sizeinsqfeet"`
	Location       pgtype.Text      `json:"location"`
	Images         []string         `json:"images"`
	Demand         pgtype.Text      `json:"demand"`
	Status         pgtype.Text      `json:"status"`
	UserID         pgtype.Int4      `json:"user_id"`
	CreatedAt      pgtype.Timestamp `json:"created_at"`
	UpdatedAt      pgtype.Timestamp `json:"updated_at"`
	ID_2           int64            `json:"id_2"`
	Name           string           `json:"name"`
	Email          pgtype.Text      `json:"email"`
	Password       pgtype.Text      `json:"password"`
	ProfilePicture pgtype.Text      `json:"profile_picture"`
	Otp            pgtype.Int4      `json:"otp"`
	CreatedAt_2    pgtype.Timestamp `json:"created_at_2"`
	UpdatedAt_2    pgtype.Timestamp `json:"updated_at_2"`
}

func (q *Queries) GetPropertyByUserID(ctx context.Context, userID pgtype.Int4) ([]GetPropertyByUserIDRow, error) {
	rows, err := q.db.Query(ctx, getPropertyByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPropertyByUserIDRow{}
	for rows.Next() {
		var i GetPropertyByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Sizeinsqfeet,
			&i.Location,
			&i.Images,
			&i.Demand,
			&i.Status,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.ProfilePicture,
			&i.Otp,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
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

const insertProperty = `-- name: InsertProperty :one
INSERT INTO properties (sizeInSqFeet, location,demand,user_id,images) VALUES ($1,$2,$3,$4,$5) RETURNING id, sizeinsqfeet, location, images, demand, status, user_id, created_at, updated_at
`

type InsertPropertyParams struct {
	Sizeinsqfeet pgtype.Int4 `json:"sizeinsqfeet"`
	Location     pgtype.Text `json:"location"`
	Demand       pgtype.Text `json:"demand"`
	UserID       pgtype.Int4 `json:"user_id"`
	Images       []string    `json:"images"`
}

func (q *Queries) InsertProperty(ctx context.Context, arg InsertPropertyParams) (Property, error) {
	row := q.db.QueryRow(ctx, insertProperty,
		arg.Sizeinsqfeet,
		arg.Location,
		arg.Demand,
		arg.UserID,
		arg.Images,
	)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Sizeinsqfeet,
		&i.Location,
		&i.Images,
		&i.Demand,
		&i.Status,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateProperty = `-- name: UpdateProperty :one
UPDATE properties SET sizeInSqFeet =$2 , location = $3 ,demand = $4 ,status=$5 , images=$6,
  updated_at = NOW()
WHERE id=$1 RETURNING id, sizeinsqfeet, location, images, demand, status, user_id, created_at, updated_at
`

type UpdatePropertyParams struct {
	ID           int64       `json:"id"`
	Sizeinsqfeet pgtype.Int4 `json:"sizeinsqfeet"`
	Location     pgtype.Text `json:"location"`
	Demand       pgtype.Text `json:"demand"`
	Status       pgtype.Text `json:"status"`
	Images       []string    `json:"images"`
}

func (q *Queries) UpdateProperty(ctx context.Context, arg UpdatePropertyParams) (Property, error) {
	row := q.db.QueryRow(ctx, updateProperty,
		arg.ID,
		arg.Sizeinsqfeet,
		arg.Location,
		arg.Demand,
		arg.Status,
		arg.Images,
	)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Sizeinsqfeet,
		&i.Location,
		&i.Images,
		&i.Demand,
		&i.Status,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePropertyDemand = `-- name: UpdatePropertyDemand :one
UPDATE properties SET demand =$2,
  updated_at = NOW()
WHERE id=$1 RETURNING id, sizeinsqfeet, location, images, demand, status, user_id, created_at, updated_at
`

type UpdatePropertyDemandParams struct {
	ID     int64       `json:"id"`
	Demand pgtype.Text `json:"demand"`
}

func (q *Queries) UpdatePropertyDemand(ctx context.Context, arg UpdatePropertyDemandParams) (Property, error) {
	row := q.db.QueryRow(ctx, updatePropertyDemand, arg.ID, arg.Demand)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Sizeinsqfeet,
		&i.Location,
		&i.Images,
		&i.Demand,
		&i.Status,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePropertyImages = `-- name: UpdatePropertyImages :one
UPDATE properties SET images =$2,
  updated_at = NOW()
WHERE id=$1 RETURNING id, sizeinsqfeet, location, images, demand, status, user_id, created_at, updated_at
`

type UpdatePropertyImagesParams struct {
	ID     int64    `json:"id"`
	Images []string `json:"images"`
}

func (q *Queries) UpdatePropertyImages(ctx context.Context, arg UpdatePropertyImagesParams) (Property, error) {
	row := q.db.QueryRow(ctx, updatePropertyImages, arg.ID, arg.Images)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Sizeinsqfeet,
		&i.Location,
		&i.Images,
		&i.Demand,
		&i.Status,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePropertyLocation = `-- name: UpdatePropertyLocation :one
UPDATE properties SET location =$2,
  updated_at = NOW()
WHERE id=$1 RETURNING id, sizeinsqfeet, location, images, demand, status, user_id, created_at, updated_at
`

type UpdatePropertyLocationParams struct {
	ID       int64       `json:"id"`
	Location pgtype.Text `json:"location"`
}

func (q *Queries) UpdatePropertyLocation(ctx context.Context, arg UpdatePropertyLocationParams) (Property, error) {
	row := q.db.QueryRow(ctx, updatePropertyLocation, arg.ID, arg.Location)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Sizeinsqfeet,
		&i.Location,
		&i.Images,
		&i.Demand,
		&i.Status,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePropertySize = `-- name: UpdatePropertySize :one
UPDATE properties SET sizeInSqFeet =$2,
  updated_at = NOW()
WHERE id=$1 RETURNING id, sizeinsqfeet, location, images, demand, status, user_id, created_at, updated_at
`

type UpdatePropertySizeParams struct {
	ID           int64       `json:"id"`
	Sizeinsqfeet pgtype.Int4 `json:"sizeinsqfeet"`
}

func (q *Queries) UpdatePropertySize(ctx context.Context, arg UpdatePropertySizeParams) (Property, error) {
	row := q.db.QueryRow(ctx, updatePropertySize, arg.ID, arg.Sizeinsqfeet)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Sizeinsqfeet,
		&i.Location,
		&i.Images,
		&i.Demand,
		&i.Status,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateStatus = `-- name: UpdateStatus :one
UPDATE properties SET status =$2  ,
  updated_at = NOW()
WHERE id=$1 RETURNING id, sizeinsqfeet, location, images, demand, status, user_id, created_at, updated_at
`

type UpdateStatusParams struct {
	ID     int64       `json:"id"`
	Status pgtype.Text `json:"status"`
}

func (q *Queries) UpdateStatus(ctx context.Context, arg UpdateStatusParams) (Property, error) {
	row := q.db.QueryRow(ctx, updateStatus, arg.ID, arg.Status)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Sizeinsqfeet,
		&i.Location,
		&i.Images,
		&i.Demand,
		&i.Status,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
