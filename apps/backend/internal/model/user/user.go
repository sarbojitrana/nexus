package user

import (
	"github.com/sarbojitrana/nexus/internal/model"
)

type User struct {
	model.Base
	ClerkID        string  `json:"clerkId" db:"clerk_id"`
	Username       string  `json:"username" db:"username"`
	DisplayName    string  `json:"displayName" db:"display_name"`
	Bio            *string `json:"bio" db:"bio"`
	AvatarKey      *string `json:"avatarKey" db:"avatar_key"`
	BannerKey      *string `json:"bannerKey" db:"banner_key"`
	FollowerCount  int     `json:"followerCount" db:"follower_count"`
	FollowingCount int     `json:"followingCount" db:"following_count"`
	PostsCount     int     `json:"postsCount" db:"posts_count"`
}

type UserSummary struct {
	Username      string `json:"username" db:"username"`
	DisplayName   string `json:"displayName" db:"display_name"`
	AvatarKey     string `json:"avatarKey" db:"avatarKey"`
	FollowerCount string `json:"followerCount" db:"follower_count"`
}
