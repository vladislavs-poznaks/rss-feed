// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
insert into posts (id,
                   title,
                   description,
                   url,
                   feed_id,
                   published_at,
                   created_at,
                   updated_at)
values ($1, $2, $3, $4, $5, $6, $7, $8) returning id, title, description, url, feed_id, published_at, created_at, updated_at
`

type CreatePostParams struct {
	ID          uuid.UUID
	Title       string
	Description sql.NullString
	Url         string
	FeedID      uuid.UUID
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.Url,
		arg.FeedID,
		arg.PublishedAt,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Url,
		&i.FeedID,
		&i.PublishedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserPosts = `-- name: GetUserPosts :many
select posts.id, posts.title, posts.description, posts.url, posts.feed_id, posts.published_at, posts.created_at, posts.updated_at
from feed_follows
         inner join posts on feed_follows.feed_id = posts.feed_id
where feed_follows.user_id = $1
order by posts.published_at desc
limit $2
`

type GetUserPostsParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetUserPosts(ctx context.Context, arg GetUserPostsParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getUserPosts, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Url,
			&i.FeedID,
			&i.PublishedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
