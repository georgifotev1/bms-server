-- +goose Up
CREATE TABLE sessions (
  session_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
  expires_at TIMESTAMP NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_user_sessions_expires_at ON sessions (expires_at);

-- +goose Down
DROP INDEX idx_user_sessions_expires_at;

DROP TABLE sessions;