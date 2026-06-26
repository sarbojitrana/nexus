package community

import (
	"time"

	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
)

type CommunityRole string

const (
	MemberRole CommunityRole = "member"
	AdminRole  CommunityRole = "admin"
)

type Community struct {
	model.Base
	AdminID     uuid.UUID `json:"adminId" db:"admin_id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description *string   `json:"description" db:"description"`
	AvatarURL   *string   `json:"avatarUrl" db:"avatar_url"`
	BannerURL   *string   `json:"bannerUrl" db:"banner_url"`
	MemberCount int       `json:"memberCount" db:"member_count"`
	PostsCount  int       `json:"postsCount" db:"posts_count"`
}

type CommunityMember struct {
	UserID      uuid.UUID `json:"userId" db:"user_id"`
	CommunityID uuid.UUID `json:"communityId" db:"community_id"`
	Role        CommunityRole    `json:"role" db:"role"`
	JoinedAt    time.Time `json:"joinedAt" db:"joined_at"`
	model.BaseWithUpdatedAt
}
