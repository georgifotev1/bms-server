-- +goose Up
CREATE TABLE notifications (
  notification_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  message TEXT NOT NULL,
  roles TEXT[],
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- +goose Down
DROP TABLE notifications;
