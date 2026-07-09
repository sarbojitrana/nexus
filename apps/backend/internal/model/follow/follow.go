package follow

import (
	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
)

type UserFollow struct {
	model.BaseWithId
	model.BaseWithCreatedAt
	FollowerID  uuid.UUID `json:"followerId" db:"follower_id"`
	FollowingID uuid.UUID `json:"followingId" db:"following_id"`
}

func (f *UserFollow) SelfFollowCheck() bool {
	return f.FollowerID != f.FollowingID
}
