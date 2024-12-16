-- name: CreateUser :exec
INSERT INTO users (first_name, last_name, email, password, role)
VALUES ($1, $2, $3, $4, $5);

-- name: GetUserById :one
SELECT * FROM users WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ConfirmUserRole :exec
UPDATE users
SET role = $2
WHERE user_id = $1;
