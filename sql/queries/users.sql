-- name: CreateUser :one
insert into users (id, name, created_at, updated_at, api_key)
values ($1, $2, $3, $4,
    encode(sha256(random()::text::bytea), 'hex')
)
returning *;

-- name: GetUserByAPIKey :one
SELECT * FROM users
WHERE api_key = $1;
