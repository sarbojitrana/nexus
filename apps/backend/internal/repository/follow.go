package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/model/follow"
	"github.com/sarbojitrana/nexus/internal/server"
)

type FollowRepository struct{
	server *server.Server
}

func NewFollowRepository (server *server.Server) *FollowRepository{
	return &FollowRepository{
		server : server,
	}
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *FollowRepository) FollowCommunity( ctx context.Context, userID uuid.UUID, payload *follow.FollowCommunityPayload ) (*follow.CommunityFollow, error) {

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
		return nil, fmt.Errorf(
			"failed to follow community %s for user_id %s: %w",
			payload.CommunityID,
			userID,
			err,
		)
	}

	follow, err := pgx.CollectExactlyOneRow( rows, pgx.RowToStructByName[follow.CommunityFollow] )

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

func (r *FollowRepository) UnFollowCommunity( ctx context.Context, userID uuid.UUID, payload *follow.UnFollowCommunityPayload) error {

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

func (r *FollowRepository) FollowUser( ctx context.Context, userID uuid.UUID, payload *follow.FollowUserPayload ) (*follow.UserFollow, error) {

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
		return nil, fmt.Errorf(
			"failed to follow user %s for user_id %s: %w",
			payload.FollowingID,
			userID,
			err,
		)
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

func (r *FollowRepository) UnFollowUser( ctx context.Context, userID uuid.UUID, payload *follow.UnFollowUserPayload) error {

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

