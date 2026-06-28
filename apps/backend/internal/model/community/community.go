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

type CommunityReportStatus string

const (
	ReportPending   CommunityReportStatus = "pending"
	ReportResolved  CommunityReportStatus = "resolved"
	ReportDismissed CommunityReportStatus = "dismissed"
)
type Community struct {
	model.Base
	AdminID      uuid.UUID `json:"adminId" db:"admin_id"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	Description  *string   `json:"description" db:"description"`
	AvatarKey    *string   `json:"avatarKey" db:"avatar_key"`
	BannerKey    *string   `json:"bannerKey" db:"banner_key"`
	MembersCount int       `json:"membersCount" db:"members_count"`
	PostsCount   int       `json:"postsCount" db:"posts_count"`
}

type CommunityMember struct {
	UserID      uuid.UUID     `json:"userId" db:"user_id"`
	CommunityID uuid.UUID     `json:"communityId" db:"community_id"`
	Role        CommunityRole `json:"role" db:"role"`
	JoinedAt    time.Time     `json:"joinedAt" db:"joined_at"`
	model.BaseWithUpdatedAt
}

type CommnunitySummary struct {
	Name         string `json:"name" db:"name"`
	AvatarKey    string `json:"avatarKey" db:"avatar_key"`
	MembersCount int    `json:"membersCount" db:"members_count"`
	PostsCount   int    `json:"postsCount" db:"posts_count"`
}

type CommnunityFollow struct {
	model.BaseWithId
	model.BaseWithCreatedAt
	FollowerID  uuid.UUID `json:"followerId" db:"follower_id"`
	CommunityID uuid.UUID `json:"communityId" db:"community_id"`
}

type CommunityReport struct {
	model.Base
	ReporterID   uuid.UUID             `json:"reporterId" db:"reporter_id"`
	CommnunityID uuid.UUID             `json:"communityId" db:"community_id"`
	PostID       uuid.UUID             `json:"postId" db:"post_id"`
	Reason       string                `json:"reason" db:"reason"`
	Status       CommunityReportStatus `json:"status" db:"status"`
}
