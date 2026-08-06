package service

import (
	"github.com/sarbojitrana/nexus/internal/lib/job"
	"github.com/sarbojitrana/nexus/internal/repository"
	"github.com/sarbojitrana/nexus/internal/server"
)

type Services struct {
	Auth         *AuthService
	User         *UserService
	Post         *PostService
	Community    *CommunityService
	Follow       *FollowService
	Chat         *ChatService
	Notification *NotificationService
	Search       *SearchService
	Storage      *StorageService
	Job          *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	notificationService := NewNotificationService(s, repos.Notification, s.Hub)
	userService := NewUserService(s, repos.User, repos.Follow, s.Search)
	postService := NewPostService(s, repos.Post, s.Search)
	communityService := NewCommunityService(s, repos.Community, s.Search)
	followService := NewFollowService(s, repos.Follow, notificationService)
	chatService := NewChatService(s, repos.Chat, repos.User, repos.Follow, s.Hub, notificationService)
	searchService := NewSearchService(s, s.Search)
	storageService := NewStorageService(s, s.Storage)

	s.Hub.OnConnect = func(userID string) { chatService.BroadcastPresence(userID, true) }
	s.Hub.OnDisconnect = func(userID string) { chatService.BroadcastPresence(userID, false) }

	return &Services{
		Job:          s.Job,
		Auth:         authService,
		User:         userService,
		Post:         postService,
		Community:    communityService,
		Follow:       followService,
		Chat:         chatService,
		Notification: notificationService,
		Search:       searchService,
		Storage:      storageService,
	}, nil
}
