package model

import (
	"time"

	"github.com/google/uuid"
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
	Data       []T    `json:"data"`
	NextCursor string `json:"nextCursor"`
	HasMore    bool   `json:"hasMore"`
}
