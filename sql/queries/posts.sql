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

-- name: GetUserPosts :many
select posts.*
from feed_follows
         inner join posts on feed_follows.feed_id = posts.feed_id
where feed_follows.user_id = $1
order by posts.published_at desc
limit $2;