-- name: CreateNotification :exec
INSERT INTO notifications (message, roles)
VALUES ($1, $2);

-- name: GetNotifications :many
SELECT *
FROM notifications
WHERE $1 = ANY(roles);

-- name: UpdateNotificationStatus :exec
UPDATE notifications
SET is_active = $2
WHERE notification_id = $1
RETURNING *;
