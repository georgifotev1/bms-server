-- +goose Up
CREATE TABLE booking_guests (
    booking_guest_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    booking_id UUID REFERENCES bookings (booking_id) ON DELETE CASCADE,
    guest_id UUID REFERENCES guests (guest_id) ON DELETE CASCADE,
    CONSTRAINT unique_guest_per_booking UNIQUE (booking_id, guest_id)
);

-- +goose Down
DROP TABLE IF EXISTS booking_guests;
