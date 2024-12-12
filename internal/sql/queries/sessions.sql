-- name: CreateSession :one
INSERT INTO sessions (session_id, user_id, expires_at)
VALUES ($1, $2, $3)
RETURNING session_id;

-- name: GetSessionById :one
SELECT * FROM sessions WHERE session_id = $1;

-- name: ClearExpiredSessions :many
DELETE FROM sessions WHERE expires_at < NOW()
RETURNING *;
