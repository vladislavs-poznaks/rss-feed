-- name: CreateFeed :one
insert into feeds (id, user_id, name, url, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetNextFetchFeeds :many
SELECT * FROM feeds
order by last_fetched_at asc nulls first
limit $1;

-- name: MarkFetchedFeed :one
update feeds set last_fetched_at = now(), updated_at = now()
where id = $1
returning *;
