package repository

import "github.com/sarbojitrana/nexus/internal/server"

type Repositories struct {
	User      *UserRepository
	Post      *PostRepository
	Community *CommunityRepository
	Follow    *FollowRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		User: NewUserRepository(s),
		Post: NewPostRepository(s),
		Community: NewCommunityRepository(s),
		Follow: NewFollowRepository(s),
	}
}
