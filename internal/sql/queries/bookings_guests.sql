-- name: CreateBookingGuest :one
INSERT INTO
    booking_guests (booking_id, guest_id)
VALUES
    ('<booking_uuid>', '<guest_uuid>') RETURNING *;

-- name: ListGuestsByBooking :many
SELECT
    g.first_name,
    g.last_name,
    g.phone
FROM
    booking_guests bg
    JOIN guests g ON bg.guest_id = g.guest_id
WHERE
    bg.booking_id = $1;

-- name: ListBookingsByGuest :many
SELECT b.booking_id, b.room_id, b.check_in_date, b.check_out_date, b.status
FROM booking_guests bg
JOIN bookings b ON bg.booking_id = b.booking_id
WHERE bg.guest_id = $1;
