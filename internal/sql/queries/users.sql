-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id;

-- name: GetUserById :one
SELECT * FROM users WHERE user_id = $1;
