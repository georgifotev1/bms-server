-- +goose Up
CREATE TABLE notifications (
  notification_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  message TEXT NOT NULL,
  roles TEXT[],
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE user_notifications (
    user_id UUID REFERENCES users(user_id),
    notification_id UUID REFERENCES notifications(notification_id),
    read_at TIMESTAMP,
    PRIMARY KEY (user_id, notification_id)
);

-- +goose Down
DROP TABLE user_notifications;
DROP TABLE notifications;
