-- +goose Up
CREATE TABLE hotels (
    hotel_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(50) NOT NULL,
    address TEXT NOT NULL
);

CREATE TABLE rooms (
    room_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    hotel_id UUID REFERENCES Hotels (hotel_id) ON DELETE CASCADE,
    room_number INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('single', 'double', 'suite')),
    price_per_night DECIMAL(10, 2) NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    CONSTRAINT unique_room_number_per_hotel UNIQUE (hotel_id, room_number)
);

CREATE TABLE guests (
    guest_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    age INTEGER NOT NULL,
    phone VARCHAR(20)
);

CREATE TABLE bookings (
    booking_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    room_id UUID REFERENCES Rooms (room_id) ON DELETE CASCADE,
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (
        status in ('pending', 'confirmed', 'completed', 'cancelled')
    )
);

-- +goose Down
DROP TABLE IF EXISTS bookings;

DROP TABLE IF EXISTS guests;

DROP TABLE IF EXISTS rooms;

DROP TABLE IF EXISTS hotels;
