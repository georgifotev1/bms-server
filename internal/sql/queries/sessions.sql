-- name: CreateSession :one
INSERT INTO sessions (session_id, user_id, expires_at)
VALUES ($1, $2, $3)
RETURNING session_id;

-- name: GetSessionById :one
SELECT * FROM sessions WHERE session_id = $1;

-- name: ClearExpiredSessions :exec
DELETE FROM sessions WHERE expires_at <= NOW() OR is_active = FALSE;

-- name: GetSessionByUserId :one
SELECT * FROM sessions WHERE user_id = $1;
