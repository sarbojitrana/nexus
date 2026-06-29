package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/post"
	"github.com/sarbojitrana/nexus/internal/model/user"
	"github.com/sarbojitrana/nexus/internal/server"
)

type UserRepository struct {
	server *server.Server
}

func NewUserRepository(server *server.Server) *UserRepository {
	return &UserRepository{
		server: server,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, payload *user.CreateUserPayload) (*user.User, error) {
	stmt := `
		INSERT INTO 
			users(
				user_name,
				display_name,
				email_id,
				clerk_id,
				bio,
				avatar_key,
				banner_key
			)
		VALUES
			(
				@user_name,
				@display_name,
				@email_id,
				@clerk_id,
				@bio,
				@avatar_key,
				@banner_key
			)
		RETURNING 
		*
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_name":    payload.Username,
		"display_name": payload.DisplayName,
		"email_id":     payload.EmailID,
		"clerk_id":     payload.ClerkID,
		"bio":          payload.Bio,
		"avatar_key":   payload.AvatarKey,
		"banner_key":   payload.BannerKey,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to create user query: %w ", err)
	}

	defer rows.Close()

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	if err != nil {
		return nil, fmt.Errorf("Failed to scan created user: %w ", err)
	}

	return &userItem, nil
}

//-------------------------------------------------------------------------------------------

func (r *UserRepository) UpdateUser(ctx context.Context, userID uuid.UUID, payload *user.UpdateUserPayload) (*user.User, error) {
	stmt := `
		UPDATE users SET 
	`
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	setClauses := []string{}

	if payload.Username != nil {
		setClauses = append(setClauses, "username = @username")
		args["username"] = payload.Username
	}

	if payload.EmailID != nil {
		setClauses = append(setClauses, "email_id = @email_id")
		args["email_id"] = payload.EmailID
	}

	if payload.DisplayName != nil {
		setClauses = append(setClauses, "display_name = @display_name")
		args["display_name"] = payload.DisplayName
	}

	if payload.Bio != nil {
		setClauses = append(setClauses, "bio = @bio")
		args["bio"] = payload.Bio
	}

	if payload.AvatarKey != nil {
		setClauses = append(setClauses, "avatar_key = @avatar_key")
		args["avatar_key"] = payload.AvatarKey
	}

	if payload.BannerKey != nil {
		setClauses = append(setClauses, "banner_key = @banner_key")
		args["banner_key"] = payload.BannerKey
	}

	if len(setClauses) == 0 {
		return nil, errs.NewBadRequestError("No field provide to be updated", false, nil, nil, nil)
	}

	stmt += strings.Join(setClauses, ", ")
	stmt += " WHERE id = @user_id RETURNING *"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to exectute the query for user_id %s : %w", userID, err)
	}

	defer rows.Close()

	updatedUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	if err != nil {
		return nil, fmt.Errorf("Failed to scan updated user for user_id %s : %w", userID, err)
	}

	return &updatedUser, nil
}

//-------------------------------------------------------------------------------------------

func (r *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	stmt := `
		SELECT *
		FROM users
		WHERE id = @user_id
	`
	args := pgx.NamedArgs{
		"user_id": userID,
	}

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute get user query for user_id %s: %w", userID, err)
	}

	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	if err != nil {
		return nil, fmt.Errorf("Failed to scan user for user_id %s: %w", userID, err)
	}

	return &user, nil
}

//-------------------------------------------------------------------------------------------

func (r *UserRepository) GetPostsByUserID(ctx context.Context, userID uuid.UUID, payload user.GetPostsByUserIDPayload) (*model.CursorPaginatedResponse[post.ViewPost], error) {
	stmt := `
		SELECT p.*,
		COALESCE(
			json_agg(
				to_json(camel(pm))
				ORDER BY pm.created_at ASC
			) FILTER(
				WHERE pm.id IS NOT NULL 
			),
			'[]'::jsonb
		) AS post_media
		FROM posts p
		LEFT JOIN post_media pm
			ON pm.post_id = p.id
	`
	limit := 20
	args := pgx.NamedArgs{
		"user_id":        userID,
		"limit_plus_one": limit + 1,
	}

	if payload.NextCursor != nil {
		args["cursor"] = *payload.NextCursor
		if *payload.Order == "asc" {
			stmt += " p.created_at >= @cursor"
		} else {
			stmt += " p.created_at <= @cursor"
		}
	}

	conditions := []string{}

	if payload.DateCreatedStart != nil {
		conditions = append(conditions, "p.created_at >= @date_created_start")
		args["date_created_start"] = *payload.DateCreatedStart
	}

	if payload.DateCreatedEnd != nil {
		conditions = append(conditions, "p.created_at <= @date_created_end")
		args["date_created_end"] = *payload.DateCreatedEnd
	}

	stmt += " WHERE p.author_id = @user_id AND p.parent_post_id IS NULL AND "
	if len(conditions) > 0 {
		stmt += strings.Join(conditions, " AND ") + " AND "
	}

	stmt += " GROUP BY p.id"

	if payload.SortBy != nil {
		stmt += " ORDER BY " + *payload.SortBy
		if payload.Order != nil {
			stmt += " " + *payload.Order
		}
	} else {
		stmt += " ORDER BY p.created_at DESC"
	}

	stmt += " LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to get posts for user id %s: %w", userID, err)
	}

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[post.ViewPost])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.CursorPaginatedResponse[post.ViewPost]{
				Data:       []post.ViewPost{},
				NextCursor: "",
				HasMore:    false,
			}, nil
		}

		return nil, fmt.Errorf("failed to collect rows from tables: todos for user_id=%s : %w", userID, err)
	}

	if len(posts) < limit+1 {
		length := len(posts)
		return &model.CursorPaginatedResponse[post.ViewPost]{
			Data:       posts[:length],
			NextCursor: "",
			HasMore:    false,
		}, nil
	}

	NextCursor := posts[20].CreatedAt.String()
	HasMore := true

	return &model.CursorPaginatedResponse[post.ViewPost]{
		Data:       posts[:20],
		NextCursor: NextCursor,
		HasMore:    HasMore,
	}, nil

}

//-------------------------------------------------------------------------------------------

func (r *UserRepository) GetUsers(ctx context.Context, payload user.GetUsersPayload) (*model.CursorPaginatedResponse[user.MiniUser], error) {

	return nil, nil
}
