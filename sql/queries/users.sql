-- name: CreateUser :one
insert into users (id, name, created_at, updated_at)
values ($1, $2, $3, $4)
returning *;
