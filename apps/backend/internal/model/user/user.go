package user

import (
	"time"

	"github.com/sarbojitrana/nexus/internal/model"
)

type ProfileVisibility string

const (
	ProfileVisibilityPublic        ProfileVisibility = "public"
	ProfileVisibilityFollowersOnly ProfileVisibility = "followers_only"
	ProfileVisibilityPrivate       ProfileVisibility = "private"
)

type GroupInvitePermission string

const (
	GroupInvitePermissionEveryone      GroupInvitePermission = "everyone"
	GroupInvitePermissionFollowersOnly GroupInvitePermission = "followers_only"
	GroupInvitePermissionNoOne         GroupInvitePermission = "no_one"
)

type User struct {
	model.BaseWithCreatedAt
	model.BaseWithUpdatedAt
	UserID                string                `json:"id" db:"id"`
	Username              string                `json:"username" db:"username"`
	EmailID               string                `json:"emailId" db:"email_id"`
	DisplayName           string                `json:"displayName" db:"display_name"`
	Bio                   *string               `json:"bio" db:"bio"`
	AvatarKey             *string               `json:"avatarKey" db:"avatar_key"`
	BannerKey             *string               `json:"bannerKey" db:"banner_key"`
	FollowerCount         int                   `json:"followerCount" db:"follower_count"`
	FollowingCount        int                   `json:"followingCount" db:"following_count"`
	PostsCount            int                   `json:"postsCount" db:"posts_count"`
	ProfileVisibility     ProfileVisibility     `json:"profileVisibility" db:"profile_visibility"`
	ShowOnlineStatus      bool                  `json:"showOnlineStatus" db:"show_online_status"`
	GroupInvitePermission GroupInvitePermission `json:"groupInvitePermission" db:"group_invite_permission"`
	ShareReadReceipts     bool                  `json:"shareReadReceipts" db:"share_read_receipts"`
}

type MiniUser struct {
	ID            string    `json:"id" db:"id"`
	Username      string    `json:"username" db:"username"`
	DisplayName   string    `json:"displayName" db:"display_name"`
	AvatarKey     string    `json:"avatarKey" db:"avatar_key"`
	Bio           string    `json:"bio" db:"bio"`
	FollowerCount int       `json:"followerCount" db:"follower_count"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

type UserResponse struct {
	model.BaseWithCreatedAt
	model.BaseWithUpdatedAt
	ID             string  `json:"id"`
	Username       string  `json:"username"`
	DisplayName    string  `json:"displayName"`
	Bio            *string `json:"bio"`
	AvatarURL      *string `json:"avatarUrl"`
	BannerURL      *string `json:"bannerUrl"`
	FollowerCount  int     `json:"followerCount" `
	FollowingCount int     `json:"followingCount"`
	PostsCount     int     `json:"postsCount"`
}

type UserBlock struct {
	BlockerId string `json:"blockerId" db:"blocker_id"`
	BlockedId string `json:"blockedId" db:"blocked_id"`
}
