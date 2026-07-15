package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/follow"
	"github.com/sarbojitrana/nexus/internal/model/user"
	"github.com/sarbojitrana/nexus/internal/server"
)

type FollowRepository struct {
	server *server.Server
}

func NewFollowRepository(server *server.Server) *FollowRepository {
	return &FollowRepository{
		server: server,
	}
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *FollowRepository) FollowCommunity(ctx context.Context, userID uuid.UUID, payload *follow.FollowCommunityPayload) (*follow.CommunityFollow, error) {

	commRepo := NewCommunityRepository(r.server)
	bannedCheck, err := commRepo.IsBannedFromCommunity(ctx, userID, payload.CommunityID)

	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was banned: %w", err)
	}

	if *bannedCheck {
		code := "USER IS BANNED"
		errs.NewBadRequestError(
			"can't follow this community",
			false,
			&code,
			nil,
			nil,
		)
	}

	stmt := `
		INSERT INTO community_follows (
			follower_id,
			community_id
		)
		VALUES (
			@follower_id,
			@community_id
		)
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"follower_id":  userID,
		"community_id": payload.CommunityID,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			code := "ALREADY_FOLLOWING"
			return nil, errs.NewBadRequestError("already following this community", false, &code, nil, nil)
		}
		return nil, fmt.Errorf("failed to follow community %s for community_id %s: %w", payload.CommunityID, userID, err)
	}

	follow, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[follow.CommunityFollow])

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse followed community for user_id %s and community_id %s: %w",
			userID,
			payload.CommunityID,
			err,
		)
	}

	return &follow, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *FollowRepository) UnFollowCommunity(ctx context.Context, userID uuid.UUID, payload *follow.UnFollowCommunityPayload) error {

	stmt := `
		DELETE FROM community_follows
		WHERE follower_id = @follower_id
		AND community_id = @community_id
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"follower_id":  userID,
		"community_id": payload.ID,
	})

	if err != nil {
		return fmt.Errorf(
			"failed to unfollow community %s for user_id %s: %w",
			payload.ID,
			userID,
			err,
		)
	}

	if result.RowsAffected() == 0 {
		code := "COMMUNITY FOLLOW NOT FOUND"
		return errs.NewNotFoundError(
			"community follow not found",
			false,
			&code,
		)
	}

	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *FollowRepository) FollowUser(ctx context.Context, userID uuid.UUID, payload *follow.FollowUserPayload) (*follow.UserFollow, error) {
	userRepo := NewUserRepository(r.server)

	blockedCheck, err := userRepo.IsUserBlocked(ctx, userID, payload.FollowingID)

	if err != nil {
		return nil, fmt.Errorf("Failed to check if the user was blocked: %w", err)
	}

	if *blockedCheck {
		code := "USER IS BLOCKED"
		errs.NewBadRequestError(
			"can't follow this user",
			false,
			&code,
			nil,
			nil,
		)
	}

	stmt := `
		INSERT INTO user_follows (
			follower_id,
			following_id
		)
		VALUES (
			@follower_id,
			@following_id
		)
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"follower_id":  userID,
		"following_id": payload.FollowingID,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			code := "ALREADY_FOLLOWING"
			return nil, errs.NewBadRequestError("already following this user", false, &code, nil, nil)
		}
		return nil, fmt.Errorf("failed to follow user %s for user_id %s: %w", payload.FollowingID, userID, err)
	}

	userFollow, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[follow.UserFollow],
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse user follow for user_id %s and following_id %s: %w",
			userID,
			payload.FollowingID,
			err,
		)
	}

	return &userFollow, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *FollowRepository) UnFollowUser(ctx context.Context, userID uuid.UUID, payload *follow.UnFollowUserPayload) error {

	stmt := `
		DELETE FROM user_follows
		WHERE follower_id = @follower_id
		AND following_id = @following_id
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"follower_id":  userID,
		"following_id": payload.FollowingID,
	})

	if err != nil {
		return fmt.Errorf(
			"failed to unfollow user %s for user_id %s: %w",
			payload.FollowingID,
			userID,
			err,
		)
	}

	if result.RowsAffected() == 0 {
		code := "USER FOLLOW NOT FOUND"
		return errs.NewNotFoundError(
			"user follow not found",
			false,
			&code,
		)
	}

	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *FollowRepository) GetFollowers(ctx context.Context, viewerID *uuid.UUID, userID uuid.UUID, query *follow.GetFollowersQuery) (*model.CursorPaginatedResponse[user.MiniUser], error) {

	stmt := `
		SELECT u.id, u.username, u.display_name, u.avatar_key, u.follower_count, u.bio, uf.created_at
		FROM user_follows uf
		JOIN users u ON u.id = uf.follower_id
		WHERE uf.following_id = @user_id
	`
	conditions := []string{}
	args := pgx.NamedArgs{"user_id": userID}

	if viewerID != nil {
		args["viewer_id"] = *viewerID
		conditions = append(conditions, `
			NOT EXISTS (
				SELECT 1 FROM user_blocks ub
				WHERE (ub.blocker_id = @viewer_id AND ub.blocked_id = u.id)
				   OR (ub.blocker_id = u.id AND ub.blocked_id = @viewer_id)
			)
		`)
	}

	orderStmt := ""
	if *query.Order == model.OrderDesc {
		orderStmt = "ORDER BY uf.created_at DESC"
		if query.CursorCreatedAt != nil {
			args["cursor_created_at"] = *query.CursorCreatedAt
			conditions = append(conditions, "uf.created_at <= @cursor_created_at")
		}
	} else {
		orderStmt = "ORDER BY uf.created_at ASC"
		if query.CursorCreatedAt != nil {
			args["cursor_created_at"] = *query.CursorCreatedAt
			conditions = append(conditions, "uf.created_at >= @cursor_created_at")
		}
	}

	if len(conditions) > 0 {
		stmt += " AND " + strings.Join(conditions, " AND ")
	}

	limit := 20
	args["limit_plus_one"] = limit + 1
	stmt += " " + orderStmt + " LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query followers for user_id %s: %w", userID, err)
	}

	followers, err := pgx.CollectRows(rows, pgx.RowToStructByName[user.MiniUser])
	if err != nil {
		return nil, fmt.Errorf("failed to collect follower rows for user_id %s: %w", userID, err)
	}

	if len(followers) < limit+1 {
		var cursorCreatedAt time.Time
		if query.CursorCreatedAt != nil {
			cursorCreatedAt = *query.CursorCreatedAt
		}
		return &model.CursorPaginatedResponse[user.MiniUser]{
			Data:            followers,
			CursorSortValue: "",
			CursorCreatedAt: cursorCreatedAt,
			HasMore:         false,
		}, nil
	}

	return &model.CursorPaginatedResponse[user.MiniUser]{
		Data:            followers[:limit],
		CursorSortValue: "",
		CursorCreatedAt: followers[limit].CreatedAt,
		HasMore:         true,
	}, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *FollowRepository) GetFollowing(ctx context.Context, viewerID *uuid.UUID, userID uuid.UUID, query *follow.GetFollowersQuery) (*model.CursorPaginatedResponse[user.MiniUser], error) {

	stmt := `
		SELECT u.id, u.username, u.display_name, u.avatar_key, u.follower_count, u.bio, uf.created_at
		FROM user_follows uf
		JOIN users u ON u.id = uf.following_id
		WHERE uf.follower_id = @user_id
	`
	conditions := []string{}
	args := pgx.NamedArgs{"user_id": userID}

	if viewerID != nil {
		args["viewer_id"] = *viewerID
		conditions = append(conditions, `
			NOT EXISTS (
				SELECT 1 FROM user_blocks ub
				WHERE (ub.blocker_id = @viewer_id AND ub.blocked_id = u.id)
				   OR (ub.blocker_id = u.id AND ub.blocked_id = @viewer_id)
			)
		`)
	}

	orderStmt := ""
	if *query.Order == model.OrderDesc {
		orderStmt = "ORDER BY uf.created_at DESC"
		if query.CursorCreatedAt != nil {
			args["cursor_created_at"] = *query.CursorCreatedAt
			conditions = append(conditions, "uf.created_at <= @cursor_created_at")
		}
	} else {
		orderStmt = "ORDER BY uf.created_at ASC"
		if query.CursorCreatedAt != nil {
			args["cursor_created_at"] = *query.CursorCreatedAt
			conditions = append(conditions, "uf.created_at >= @cursor_created_at")
		}
	}

	if len(conditions) > 0 {
		stmt += " AND " + strings.Join(conditions, " AND ")
	}

	limit := 20
	args["limit_plus_one"] = limit + 1
	stmt += " " + orderStmt + " LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query following for user_id %s: %w", userID, err)
	}

	following, err := pgx.CollectRows(rows, pgx.RowToStructByName[user.MiniUser])
	if err != nil {
		return nil, fmt.Errorf("failed to collect following rows for user_id %s: %w", userID, err)
	}

	if len(following) < limit+1 {
		var cursorCreatedAt time.Time
		if query.CursorCreatedAt != nil {
			cursorCreatedAt = *query.CursorCreatedAt
		}
		return &model.CursorPaginatedResponse[user.MiniUser]{
			Data:            following,
			CursorSortValue: "",
			CursorCreatedAt: cursorCreatedAt,
			HasMore:         false,
		}, nil
	}

	return &model.CursorPaginatedResponse[user.MiniUser]{
		Data:            following[:limit],
		CursorSortValue: "",
		CursorCreatedAt: following[limit].CreatedAt,
		HasMore:         true,
	}, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
