-- name: CreatePost :one
insert into posts (id,
                   title,
                   description,
                   url,
                   feed_id,
                   published_at,
                   created_at,
                   updated_at)
values ($1, $2, $3, $4, $5, $6, $7, $8) returning *;
