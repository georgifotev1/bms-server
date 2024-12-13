-- name: CreateNotification :exec
INSERT INTO notifications (message, roles)
VALUES ($1, $2);

-- name: GetNotificationsByUserId :many
SELECT
    n.notification_id,
    n.message
FROM notifications n
LEFT JOIN user_notifications un ON
    n.notification_id = un.notification_id AND
    un.user_id = $1
WHERE
    n.is_active = true AND
    (un.read_at IS NULL);

-- name: MarkAsRead :exec
INSERT INTO user_notifications (user_id, notification_id, read_at)
VALUES ($1, $2, NOW())
ON CONFLICT (user_id, notification_id)
DO UPDATE SET read_at = NOW();
