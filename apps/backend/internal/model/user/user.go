package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
)

type User struct {
	model.Base
	ClerkID        string  `json:"clerkId" db:"clerk_id"`
	Username       string  `json:"username" db:"username"`
	EmailID        string  `json:"emailId" db:"email_id"`
	DisplayName    string  `json:"displayName" db:"display_name"`
	Bio            *string `json:"bio" db:"bio"`
	AvatarKey      *string `json:"avatarKey" db:"avatar_key"`
	BannerKey      *string `json:"bannerKey" db:"banner_key"`
	FollowerCount  int     `json:"followerCount" db:"follower_count"`
	FollowingCount int     `json:"followingCount" db:"following_count"`
	PostsCount     int     `json:"postsCount" db:"posts_count"`
}

type MiniUser struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	Username      string    `json:"username" db:"username"`
	DisplayName   string    `json:"displayName" db:"display_name"`
	AvatarKey     string    `json:"avatarKey" db:"avatar_key"`
	Bio           string    `json:"bio" db:"bio"`
	FollowerCount int       `json:"followerCount" db:"follower_count"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

type UserResponse struct {
	model.Base
	Username       string  `json:"username"`
	DisplayName    string  `json:"displayName"`
	Bio            *string `json:"bio"`
	AvatarURL      *string `json:"avatarUrl"`
	BannerURL      *string `json:"bannerUrl"`
	FollowerCount  int     `json:"followerCount" `
	FollowingCount int     `json:"followingCount"`
	PostsCount     int     `json:"postsCount"`
}

type UserBlock struct{
	BlockerId 	uuid.UUID 	`json:"blockerId" db:"blocker_id"`
	BlockedId 	uuid.UUID 	`json:"blockedId" db:"blocked_id"`
}
