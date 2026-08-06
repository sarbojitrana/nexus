package handler

import (
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
)

type Handlers struct {
	Health       *HealthHandler
	OpenAPI      *OpenAPIHandler
	Webhook      *ClerkWebhookHandler
	User         *UserHandler
	Post         *PostHandler
	Community    *CommunityHandler
	Follow       *FollowHandler
	Chat         *ChatHandler
	Notification *NotificationHandler
	Search       *SearchHandler
	Storage      *StorageHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:       NewHealthHandler(s),
		OpenAPI:      NewOpenAPIHandler(s),
		Webhook:      NewClerkWebhookHandler(s, services),
		User:         NewUserHandler(s, services),
		Post:         NewPostHandler(s, services),
		Community:    NewCommunityHandler(s, services),
		Follow:       NewFollowHandler(s, services),
		Chat:         NewChatHandler(s, services),
		Notification: NewNotificationHandler(s, services),
		Search:       NewSearchHandler(s, services),
		Storage:      NewStorageHandler(s, services),
	}
}
