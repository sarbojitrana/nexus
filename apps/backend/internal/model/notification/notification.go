package notification

import (
	"encoding/json"

	"github.com/sarbojitrana/nexus/internal/model"
)

type Type string

const (
	TypeFollow          Type = "follow"
	TypeGroupInvitation Type = "group_invitation"
	TypeMessage         Type = "message"
)

type Notification struct {
	model.BaseWithId
	UserID  string          `json:"userId" db:"user_id"`
	ActorID *string         `json:"actorId" db:"actor_id"`
	Type    Type            `json:"type" db:"type"`
	Data    json.RawMessage `json:"data" db:"data"`
	IsRead  bool            `json:"isRead" db:"is_read"`
	model.BaseWithCreatedAt
}
