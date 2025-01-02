-- name: CreateSession :one
INSERT INTO sessions (session_id, user_id, expires_at)
VALUES ($1, $2, $3)
RETURNING session_id;

-- name: GetSessionById :one
SELECT * FROM sessions WHERE session_id = $1;

-- name: GetSessionByUserId :one
SELECT * FROM sessions WHERE user_id = $1;

-- name: UpdateSession :one
UPDATE sessions
SET expires_at = $2
WHERE session_id = $1
RETURNING session_id;
