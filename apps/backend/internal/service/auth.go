package service

import (
	"context"
	"fmt"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/server"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
)

type AuthService struct {
	server *server.Server
}

func NewAuthService(s *server.Server) *AuthService {
	clerk.SetKey(s.Config.Auth.SecretKey)
	return &AuthService{
		server: s,
	}
}

func (s *AuthService) GetUserEmail(ctx context.Context, userID string) (string, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	user, err := clerkUser.Get(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Str("clerk_user_id", userID).Msg("failed to fetch user from clerk")
		return "", fmt.Errorf("Coudld not get user email address from Clerk: %w", err)
	}

	return primaryEmailFromClerkUser(user)
}
