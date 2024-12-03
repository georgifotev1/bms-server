-- name: CreateGuest :one
INSERT INTO guests (first_name, last_name, age, phone)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetGuestById :one
SELECT * FROM guests WHERE guest_id = $1;

-- name: ListGuests :many
SELECT * FROM guests LIMIT $1 OFFSET $2;

-- name: UpdateGuest :one
UPDATE guests
SET phone = $2
WHERE guest_id = $1
RETURNING *;

-- name: DeleteGuest :one
DELETE FROM guests WHERE guest_id = $1
RETURNING guest_id;
