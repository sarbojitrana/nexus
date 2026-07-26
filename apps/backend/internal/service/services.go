package service

import (
	"github.com/sarbojitrana/nexus/internal/lib/job"
	"github.com/sarbojitrana/nexus/internal/repository"
	"github.com/sarbojitrana/nexus/internal/server"
)

type Services struct {
	Auth      *AuthService
	User      *UserService
	Post      *PostService
	Community *CommunityService
	Follow    *FollowService
	Job       *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	userService := NewUserService(s, repos.User)
	postService := NewPostService(s, repos.Post)
	communityService := NewCommunityService(s, repos.Community)
	followService := NewFollowService(s, repos.Follow)

	return &Services{
		Job:       s.Job,
		Auth:      authService,
		User:      userService,
		Post:      postService,
		Community: communityService,
		Follow:    followService,
	}, nil
}
