package notification

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type GetNotificationsQuery struct {
	CursorCreatedAt *time.Time `query:"cursorCreatedAt"`
}

func (p *GetNotificationsQuery) Validate() error {
	return validator.New().Struct(p)
}
