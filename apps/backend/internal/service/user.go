package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

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
	"github.com/sarbojitrana/nexus/internal/search"
	"github.com/sarbojitrana/nexus/internal/server"
)

const signInTimeFormat = "Jan 2, 2006, 3:04 PM MST"

type UserService struct {
	server     *server.Server
	repo       *repository.UserRepository
	followRepo *repository.FollowRepository
	search     *search.Client
}

func NewUserService(s *server.Server, repo *repository.UserRepository, followRepo *repository.FollowRepository, search *search.Client) *UserService {
	return &UserService{
		server:     s,
		repo:       repo,
		followRepo: followRepo,
		search:     search,
	}
}

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

	s.indexUser(ctx, created)
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

func (s *UserService) NotifySignIn(ctx context.Context, userID string) {
	logger := middleware.GetLoggerFromContext(ctx)

	u, err := s.repo.GetUserByID(ctx, nil, userID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to look up user for sign-in email")
		return
	}

	task, err := job.NewSignInEmailTask(u.EmailID, u.DisplayName, time.Now().UTC().Format(signInTimeFormat))
	if err != nil {
		logger.Error().Err(err).Str("to", u.EmailID).Msg("failed to build sign-in email task")
		return
	}

	if _, err := s.server.Job.Client.Enqueue(task); err != nil {
		logger.Error().Err(err).Str("to", u.EmailID).Msg("failed to enqueue sign-in email task")
		return
	}

	logger.Info().Str("event", "signin_email_enqueued").Str("user_id", userID).Msg("sign-in email enqueued")
}

func (s *UserService) NotifyPasswordChanged(ctx context.Context, userID string) error {
	logger := middleware.GetLoggerFromContext(ctx)

	u, err := s.repo.GetUserByID(ctx, nil, userID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to look up user for password-changed email")
		return err
	}

	task, err := job.NewPasswordChangedEmailTask(u.EmailID, u.DisplayName, time.Now().UTC().Format(signInTimeFormat))
	if err != nil {
		logger.Error().Err(err).Str("to", u.EmailID).Msg("failed to build password-changed email task")
		return err
	}

	if _, err := s.server.Job.Client.Enqueue(task); err != nil {
		logger.Error().Err(err).Str("to", u.EmailID).Msg("failed to enqueue password-changed email task")
		return err
	}

	logger.Info().Str("event", "password_changed_email_enqueued").Str("user_id", userID).Msg("password-changed email enqueued")
	return nil
}

func (s *UserService) indexUser(ctx context.Context, u *user.User) {
	if s.search == nil {
		return
	}
	doc := search.UserDoc{
		ID:            u.UserID,
		Username:      u.Username,
		DisplayName:   u.DisplayName,
		Bio:           u.Bio,
		FollowerCount: u.FollowerCount,
		CreatedAt:     u.CreatedAt,
	}
	if err := s.search.IndexUser(ctx, doc); err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("user_id", u.UserID).Msg("failed to index user")
	}
}

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

	if s.search != nil {
		if err := s.search.DeleteUser(ctx, clerkID); err != nil {
			logger.Error().Err(err).Str("user_id", clerkID).Msg("failed to remove user from search index")
		}
	}

	return nil
}

func (s *UserService) GetByID(ctx context.Context, viewerID *string, userID string) (*user.User, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	u, err := s.repo.GetUserByID(ctx, viewerID, userID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to fetch user")
		return nil, err
	}

	visible, err := s.canViewProfile(ctx, viewerID, u)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to check profile visibility")
		return nil, err
	}
	if !visible {
		code := "USER_NOT_FOUND"
		return nil, errs.NewNotFoundError("user not found", false, &code)
	}

	return u, nil
}

func (s *UserService) canViewProfile(ctx context.Context, viewerID *string, target *user.User) (bool, error) {
	if viewerID != nil && *viewerID == target.UserID {
		return true, nil
	}

	switch target.ProfileVisibility {
	case user.ProfileVisibilityPrivate:
		return false, nil
	case user.ProfileVisibilityFollowersOnly:
		if viewerID == nil {
			return false, nil
		}
		isFollowing, err := s.followRepo.IsFollowingUser(ctx, *viewerID, target.UserID)
		if err != nil {
			return false, err
		}
		return *isFollowing, nil
	default: // public
		return true, nil
	}
}

func (s *UserService) List(ctx context.Context, viewerID *string, query *user.GetUsersQuery) (*model.CursorPaginatedResponse[user.MiniUser], error) {
	res, err := s.repo.GetUsers(ctx, viewerID, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Msg("failed to list users")
		return nil, err
	}
	return res, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, payload *user.UpdateUserPayload) (*user.User, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	updated, err := s.repo.UpdateUser(ctx, userID, payload)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to update profile")
		return nil, err
	}

	logger.Info().Str("event", "user_profile_updated").Str("user_id", userID).Msg("profile updated")
	s.indexUser(ctx, updated)
	return updated, nil
}

func (s *UserService) GetSettings(ctx context.Context, userID string) (*user.User, error) {
	u, err := s.repo.GetUserByID(ctx, nil, userID)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("user_id", userID).Msg("failed to fetch settings")
		return nil, err
	}
	return u, nil
}

func (s *UserService) UpdateSettings(ctx context.Context, userID string, payload *user.UpdateUserSettingsPayload) (*user.User, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	updated, err := s.repo.UpdateUserSettings(ctx, userID, payload)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to update settings")
		return nil, err
	}

	logger.Info().Str("event", "user_settings_updated").Str("user_id", userID).Msg("settings updated")
	return updated, nil
}

func (s *UserService) GetPostsByUser(ctx context.Context, viewerID *string, profileUserID string, payload *user.GetPostsByUserIDPayload) (*model.CursorPaginatedResponse[post.PopulatedPost], error) {
	res, err := s.repo.GetPostsByUserID(ctx, viewerID, profileUserID, payload)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("user_id", profileUserID).Msg("failed to fetch user's posts")
		return nil, err
	}
	return res, nil
}

func (s *UserService) DeleteAccount(ctx context.Context, userID string) error {
	logger := middleware.GetLoggerFromContext(ctx)

	if _, err := clerkUser.Delete(ctx, userID); err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to delete clerk account")
		return fmt.Errorf("failed to delete clerk account for user %s: %w", userID, err)
	}

	logger.Info().Str("event", "user_account_deletion_requested").Str("user_id", userID).Msg("clerk account deletion requested")
	return nil
}
