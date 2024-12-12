-- +goose Up
CREATE TABLE notifications (
  notification_id UUID PRIMARY KEY NOT NULL,
  message TEXT NOT NULL,
  roles TEXT[],
  user_id UUID REFERENCES users (user_id) ON DELETE CASCADE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- +goose Down
DROP TABLE notifications;
