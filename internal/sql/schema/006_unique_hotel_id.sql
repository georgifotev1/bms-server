-- +goose Up
ALTER TABLE rooms
ALTER COLUMN hotel_id SET NOT NULL;

-- +goose Down
ALTER TABLE rooms
ALTER COLUMN hotel_id DROP NOT NULL;
