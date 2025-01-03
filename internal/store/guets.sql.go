// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: guets.sql

package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createGuest = `-- name: CreateGuest :one
INSERT INTO guests (first_name, last_name, age, phone)
VALUES ($1, $2, $3, $4)
RETURNING guest_id, first_name, last_name, age, phone
`

type CreateGuestParams struct {
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Age       int32          `json:"age"`
	Phone     sql.NullString `json:"phone"`
}

func (q *Queries) CreateGuest(ctx context.Context, arg CreateGuestParams) (Guest, error) {
	row := q.db.QueryRowContext(ctx, createGuest,
		arg.FirstName,
		arg.LastName,
		arg.Age,
		arg.Phone,
	)
	var i Guest
	err := row.Scan(
		&i.GuestID,
		&i.FirstName,
		&i.LastName,
		&i.Age,
		&i.Phone,
	)
	return i, err
}

const deleteGuest = `-- name: DeleteGuest :one
DELETE FROM guests WHERE guest_id = $1
RETURNING guest_id
`

func (q *Queries) DeleteGuest(ctx context.Context, guestID uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, deleteGuest, guestID)
	var guest_id uuid.UUID
	err := row.Scan(&guest_id)
	return guest_id, err
}

const getGuestById = `-- name: GetGuestById :one
SELECT guest_id, first_name, last_name, age, phone FROM guests WHERE guest_id = $1
`

func (q *Queries) GetGuestById(ctx context.Context, guestID uuid.UUID) (Guest, error) {
	row := q.db.QueryRowContext(ctx, getGuestById, guestID)
	var i Guest
	err := row.Scan(
		&i.GuestID,
		&i.FirstName,
		&i.LastName,
		&i.Age,
		&i.Phone,
	)
	return i, err
}

const listGuests = `-- name: ListGuests :many
SELECT guest_id, first_name, last_name, age, phone FROM guests LIMIT $1 OFFSET $2
`

type ListGuestsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListGuests(ctx context.Context, arg ListGuestsParams) ([]Guest, error) {
	rows, err := q.db.QueryContext(ctx, listGuests, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Guest
	for rows.Next() {
		var i Guest
		if err := rows.Scan(
			&i.GuestID,
			&i.FirstName,
			&i.LastName,
			&i.Age,
			&i.Phone,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateGuest = `-- name: UpdateGuest :one
UPDATE guests
SET phone = $2
WHERE guest_id = $1
RETURNING guest_id, first_name, last_name, age, phone
`

type UpdateGuestParams struct {
	GuestID uuid.UUID      `json:"guest_id"`
	Phone   sql.NullString `json:"phone"`
}

func (q *Queries) UpdateGuest(ctx context.Context, arg UpdateGuestParams) (Guest, error) {
	row := q.db.QueryRowContext(ctx, updateGuest, arg.GuestID, arg.Phone)
	var i Guest
	err := row.Scan(
		&i.GuestID,
		&i.FirstName,
		&i.LastName,
		&i.Age,
		&i.Phone,
	)
	return i, err
}
