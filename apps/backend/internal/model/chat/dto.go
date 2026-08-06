package chat

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type StartDirectConversationPayload struct {
	UserID string `json:"userId" validate:"required"`
}

func (p *StartDirectConversationPayload) Validate() error {
	return validator.New().Struct(p)
}

type CreateGroupConversationPayload struct {
	Name       string   `json:"name" validate:"required,max=100"`
	InviteeIDs []string `json:"inviteeIds" validate:"required,min=1,max=50,dive,required"`
}

func (p *CreateGroupConversationPayload) Validate() error {
	return validator.New().Struct(p)
}

type InviteToConversationPayload struct {
	UserID string `json:"userId" validate:"required"`
}

func (p *InviteToConversationPayload) Validate() error {
	return validator.New().Struct(p)
}

type SendMessagePayload struct {
	Content            *string `json:"content" validate:"omitempty,max=4000"`
	AttachmentKey      *string `json:"attachmentKey"`
	AttachmentMimeType *string `json:"attachmentMimeType"`
	AttachmentFileSize *int64  `json:"attachmentFileSize"`
}

func (p *SendMessagePayload) Validate() error {
	if err := validator.New().Struct(p); err != nil {
		return err
	}
	if p.Content == nil && p.AttachmentKey == nil {
		return fmt.Errorf("a message needs content, an attachment, or both")
	}
	return nil
}

type GetMessagesQuery struct {
	CursorCreatedAt *time.Time `query:"cursorCreatedAt"`
}

func (p *GetMessagesQuery) Validate() error {
	return validator.New().Struct(p)
}
