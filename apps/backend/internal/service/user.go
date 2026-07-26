package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/jackc/pgx/v5"
	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/lib/job"
	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/post"
	"github.com/sarbojitrana/nexus/internal/model/user"
	"github.com/sarbojitrana/nexus/internal/repository"
	"github.com/sarbojitrana/nexus/internal/server"
)

type UserService struct {
	server *server.Server
	repo   *repository.UserRepository
}

func NewUserService(s *server.Server, repo *repository.UserRepository) *UserService {
	return &UserService{
		server: s,
		repo:   repo,
	}
}

// CreateFromClerk provisions a minimal local user row from a Clerk user.created
// event. It's idempotent so Clerk's at-least-once webhook delivery can't create
// duplicate rows.
func (s *UserService) CreateFromClerk(ctx context.Context, cu *clerk.User) (*user.User, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	if existing, err := s.repo.GetUserByID(ctx, nil, cu.ID); err == nil {
		return existing, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		logger.Error().Err(err).Str("clerk_user_id", cu.ID).Msg("failed to check for existing user")
		return nil, err
	}

	email, err := primaryEmailFromClerkUser(cu)
	if err != nil {
		logger.Error().Err(err).Str("clerk_user_id", cu.ID).Msg("failed to extract primary email from clerk user")
		return nil, err
	}

	displayName := displayNameFromClerkUser(cu, email)

	created, err := s.repo.CreateUser(ctx, &user.CreateUserPayload{
		ClerkID:     cu.ID,
		Username:    generatePlaceholderUsername(cu.ID),
		DisplayName: displayName,
		EmailID:     email,
	})
	if err != nil {
		logger.Error().Err(err).Str("clerk_user_id", cu.ID).Msg("failed to create user")
		return nil, err
	}

	logger.Info().Str("event", "user_created").Str("user_id", created.UserID).Msg("user created from clerk webhook")

	s.enqueueWelcomeEmail(ctx, email, displayName)

	return created, nil
}

func (s *UserService) enqueueWelcomeEmail(ctx context.Context, email, firstName string) {
	logger := middleware.GetLoggerFromContext(ctx)

	task, err := job.NewWelcomeEmailTask(email, firstName)
	if err != nil {
		logger.Error().Err(err).Str("to", email).Msg("failed to build welcome email task")
		return
	}

	if _, err := s.server.Job.Client.Enqueue(task); err != nil {
		logger.Error().Err(err).Str("to", email).Msg("failed to enqueue welcome email task")
	}
}

// UpdateFromClerk syncs the local email copy from a Clerk user.updated event.
// Email is the only field Clerk stays the source of truth for post-onboarding.
func (s *UserService) UpdateFromClerk(ctx context.Context, cu *clerk.User) (*user.User, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	email, err := primaryEmailFromClerkUser(cu)
	if err != nil {
		logger.Error().Err(err).Str("clerk_user_id", cu.ID).Msg("failed to extract primary email from clerk user")
		return nil, err
	}

	updated, err := s.repo.UpdateUserEmail(ctx, cu.ID, email)
	if err != nil {
		logger.Error().Err(err).Str("clerk_user_id", cu.ID).Msg("failed to sync user email")
		return nil, err
	}

	logger.Info().Str("event", "user_email_synced").Str("user_id", cu.ID).Msg("user email synced from clerk webhook")
	return updated, nil
}

// DeleteFromClerk removes the local user row for a Clerk user.deleted event.
// A missing row is treated as a no-op so redelivered events stay idempotent.
func (s *UserService) DeleteFromClerk(ctx context.Context, clerkID string) error {
	logger := middleware.GetLoggerFromContext(ctx)

	err := s.repo.DeleteUser(ctx, clerkID)

	var httpErr *errs.HTTPError
	if errors.As(err, &httpErr) && httpErr.Status == http.StatusNotFound {
		return nil
	}
	if err != nil {
		logger.Error().Err(err).Str("clerk_user_id", clerkID).Msg("failed to delete user")
		return err
	}

	logger.Info().Str("event", "user_deleted").Str("user_id", clerkID).Msg("user deleted from clerk webhook")
	return nil
}

// GetByID fetches a profile. viewerID, when present, excludes profiles that
// have blocked (or been blocked by) the viewer.
func (s *UserService) GetByID(ctx context.Context, viewerID *string, userID string) (*user.User, error) {
	u, err := s.repo.GetUserByID(ctx, viewerID, userID)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("user_id", userID).Msg("failed to fetch user")
		return nil, err
	}
	return u, nil
}

// List searches/browses user profiles.
func (s *UserService) List(ctx context.Context, viewerID *string, query *user.GetUsersQuery) (*model.CursorPaginatedResponse[user.MiniUser], error) {
	res, err := s.repo.GetUsers(ctx, viewerID, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Msg("failed to list users")
		return nil, err
	}
	return res, nil
}

// UpdateProfile applies a profile edit for the authenticated user. Email is
// deliberately not part of UpdateUserPayload -- see UpdateFromClerk.
func (s *UserService) UpdateProfile(ctx context.Context, userID string, payload *user.UpdateUserPayload) (*user.User, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	updated, err := s.repo.UpdateUser(ctx, userID, payload)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to update profile")
		return nil, err
	}

	logger.Info().Str("event", "user_profile_updated").Str("user_id", userID).Msg("profile updated")
	return updated, nil
}

// GetPostsByUser lists a profile's posts.
func (s *UserService) GetPostsByUser(ctx context.Context, viewerID *string, profileUserID string, payload *user.GetPostsByUserIDPayload) (*model.CursorPaginatedResponse[post.PopulatedPost], error) {
	res, err := s.repo.GetPostsByUserID(ctx, viewerID, profileUserID, payload)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("user_id", profileUserID).Msg("failed to fetch user's posts")
		return nil, err
	}
	return res, nil
}

// DeleteAccount deletes the user's Clerk account. The local row is removed
// asynchronously by the user.deleted webhook (DeleteFromClerk) -- Clerk stays
// the single source of truth, there's no separate local-delete code path.
func (s *UserService) DeleteAccount(ctx context.Context, userID string) error {
	logger := middleware.GetLoggerFromContext(ctx)

	if _, err := clerkUser.Delete(ctx, userID); err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to delete clerk account")
		return fmt.Errorf("failed to delete clerk account for user %s: %w", userID, err)
	}

	logger.Info().Str("event", "user_account_deletion_requested").Str("user_id", userID).Msg("clerk account deletion requested")
	return nil
}
