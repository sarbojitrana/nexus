package model

import (
	"time"

	"github.com/google/uuid"
)

type Sort string

const (
	SortByCreatedAt Sort = "created_at"
	SortByUpvotes         Sort = "upvotes"
	SortByFollowerCount   Sort = "follower_count"
	SortByFollowingCount  Sort = "following_count"
	SortByMembersCount    Sort = "members_count"
	SortByPostsCount      Sort = "posts_count"
	SortByPopularity      Sort = "popularity"
)

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

type BaseWithId struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type BaseWithCreatedAt struct {
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type BaseWithUpdatedAt struct {
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Base struct {
	BaseWithId
	BaseWithCreatedAt
	BaseWithUpdatedAt
}

type OffsetPaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type CursorPaginatedResponse[T any] struct {
	Data            []T        `json:"data"`
	CursorSortValue string    `json:"cursorSortValue"`
	CursorCreatedAt time.Time `json:"createdAt"`
	HasMore         bool       `json:"hasMore"`
}
