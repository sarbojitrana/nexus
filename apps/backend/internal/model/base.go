package model

import (
	"time"

	"github.com/google/uuid"
)

type SortBy string

const (
	SortByCreatedAt      SortBy = "created_at"
	SortByUpvotes        SortBy = "upvotes"
	SortByFollowerCount  SortBy = "follower_count"
	SortByFollowingCount SortBy = "following_count"
	SortByMembersCount   SortBy = "members_count"
	SortByPostsCount     SortBy = "posts_count"
)

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

type Cursor struct {
	SortValue string    `json:"sortValue"`
	CreatedAt time.Time `json:"createdAt"`
}

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
	Data []T `json:"data"`
	Cursor
	HasMore bool `json:"hasMore"`
}
