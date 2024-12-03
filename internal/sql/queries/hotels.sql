-- name: CreateHotel :one
INSERT INTO hotels (name, address)
VALUES ($1, $2)
RETURNING *;

-- name: GetHotelById :one
SELECT * FROM hotels WHERE hotel_id = $1;

-- name: ListHotels :many
SELECT * FROM hotels LIMIT $1 OFFSET $2;

-- name: UpdateHotel :one
UPDATE hotels
SET name = $2, address = $3
WHERE hotel_id = $1
RETURNING *;

-- name: DeleteHotel :one
DELETE FROM hotels WHERE hotel_id = $1
RETURNING hotel_id;
