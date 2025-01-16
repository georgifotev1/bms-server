-- name: CreateRoom :one
INSERT INTO rooms (hotel_id, room_number, type, price_per_night, is_available)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRoomById :one
SELECT * FROM rooms WHERE room_id = $1;

-- name: ListRooms :many
SELECT * FROM rooms WHERE hotel_id = $1 ORDER BY room_number ASC LIMIT $2 OFFSET $3;

-- name: UpdateRoom :one
UPDATE rooms
SET is_available = $2
WHERE room_id = $1
RETURNING *;

-- name: DeleteRoom :one
DELETE FROM rooms WHERE room_id = $1
RETURNING room_id;
