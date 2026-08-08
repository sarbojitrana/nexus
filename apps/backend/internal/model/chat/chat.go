package chat

import (
	"time"

	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
)

type Conversation struct {
	model.BaseWithId
	IsGroup   bool    `json:"isGroup" db:"is_group"`
	Name      *string `json:"name" db:"name"`
	CreatedBy string  `json:"createdBy" db:"created_by"`
	model.BaseWithCreatedAt
	model.BaseWithUpdatedAt
}

type Participant struct {
	ConversationID uuid.UUID  `json:"conversationId" db:"conversation_id"`
	UserID         string     `json:"userId" db:"user_id"`
	JoinedAt       time.Time  `json:"joinedAt" db:"joined_at"`
	LastReadAt     *time.Time `json:"lastReadAt" db:"last_read_at"`
}

type InvitationStatus string

const (
	InvitationPending  InvitationStatus = "pending"
	InvitationAccepted InvitationStatus = "accepted"
	InvitationRejected InvitationStatus = "rejected"
)

type Invitation struct {
	model.BaseWithId
	ConversationID uuid.UUID        `json:"conversationId" db:"conversation_id"`
	InviterID      string           `json:"inviterId" db:"inviter_id"`
	InviteeID      string           `json:"inviteeId" db:"invitee_id"`
	Status         InvitationStatus `json:"status" db:"status"`
	model.BaseWithCreatedAt
	RespondedAt *time.Time `json:"respondedAt" db:"responded_at"`
}

type Message struct {
	model.BaseWithId
	ConversationID     uuid.UUID `json:"conversationId" db:"conversation_id"`
	SenderID           string    `json:"senderId" db:"sender_id"`
	Content            *string   `json:"content" db:"content"`
	AttachmentKey      *string   `json:"attachmentKey" db:"attachment_key"`
	AttachmentMimeType *string   `json:"attachmentMimeType" db:"attachment_mime_type"`
	AttachmentFileSize *int64    `json:"attachmentFileSize" db:"attachment_file_size"`
	model.BaseWithCreatedAt
}

type ParticipantView struct {
	UserID      string     `json:"userId"`
	Username    string     `json:"username"`
	DisplayName string     `json:"displayName"`
	AvatarKey   *string    `json:"avatarKey"`
	AvatarURL   *string    `json:"avatarUrl"`
	IsOnline    bool       `json:"isOnline"`
	LastReadAt  *time.Time `json:"lastReadAt,omitempty"`
}

type ConversationSummary struct {
	Conversation
	LastMessage  *Message          `json:"lastMessage"`
	UnreadCount  int               `json:"unreadCount"`
	Participants []ParticipantView `json:"participants"`
}
