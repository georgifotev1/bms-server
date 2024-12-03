-- name: CreateBooking :one
INSERT INTO bookings (room_id, check_in_date, check_out_date, total_price, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetBookingById :one
SELECT * FROM bookings WHERE booking_id = $1;

-- name: ListBookingsByHotel :many
SELECT
    b.booking_id,
    b.room_id,
    b.check_in_date,
    b.check_out_date,
    b.total_price,
    b.status,
    r.room_number,
    r.type AS room_type,
    h.hotel_id,
    h.name AS hotel_name,
    h.address AS hotel_address
FROM bookings b
JOIN rooms r ON b.room_id = r.room_id
JOIN hotels h ON r.hotel_id = h.hotel_id
WHERE h.hotel_id = $1;

-- name: UpdateBooking :one
UPDATE bookings
SET status = $2
WHERE room_id = $1
RETURNING *;

-- name: DeleteBooking :one
DELETE FROM bookings WHERE booking_id = $1
RETURNING booking_id;
