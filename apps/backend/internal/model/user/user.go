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
	AvatarURL      *string `json:"avatarUrl" db:"avatar_url"`
	BannerURL      *string `json:"bannerUrl" db:"banner_url"`
	FollowerCount  int     `json:"followerCount" db:"follower_count"`
	FollowingCount int     `json:"followingCount" db:"following_count"`
	PostsCount     int     `json:"postsCount" db:"posts_count"`
}
