-- +goose Up
CREATE TABLE users (
  user_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  first_name VARCHAR(50) NOT NULL,
  last_name VARCHAR(50) NOT NULL,
  email VARCHAR(50) NOT NULL UNIQUE,
  password BYTEA NOT NULL,
  role VARCHAR(50) NOT NULL CHECK (role IN ('manager', 'admin', 'staff', 'pending')),
  created_at TIMESTAMP(0) NOT NULL DEFAULT NOW ()
);

-- +goose Down
DROP TABLE users;
