// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, feed_id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	FeedID    uuid.UUID
	UserID    uuid.UUID
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.FeedID,
		arg.UserID,
	)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const deleteFeedFollow = `-- name: DeleteFeedFollow :one
DELETE FROM feed_follows
WHERE id = $1
RETURNING id, created_at, updated_at, user_id, feed_id
`

func (q *Queries) DeleteFeedFollow(ctx context.Context, id uuid.UUID) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, deleteFeedFollow, id)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const getFeedFollowsByUser = `-- name: GetFeedFollowsByUser :many
SELECT id, created_at, updated_at, user_id, feed_id FROM feed_follows
WHERE user_id = $1
`

func (q *Queries) GetFeedFollowsByUser(ctx context.Context, userID uuid.UUID) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
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
