package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *UserRepository) GetPostsByUserID(ctx context.Context, userID uuid.UUID, payload *user.GetPostsByUserIDPayload) (*model.CursorPaginatedResponse[post.PopulatedPost], error) {
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
	limit := 10
	args := pgx.NamedArgs{
		"user_id":        userID,
		"limit_plus_one": limit + 1,
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

	stmt += " WHERE p.author_id = @user_id AND p.parent_post_id IS NULL"

	if len(conditions) > 0 {
		stmt += " AND " + strings.Join(conditions, " AND ")
	}

	orderStmt := ""

	if payload.Sort == nil || *payload.Sort == model.SortByCreatedAt {
		if payload.Order == nil || *payload.Order == model.OrderDesc {
			orderStmt += " ORDER BY p.created_at DESC "
		} else {
			orderStmt += " ORDER BY p.created_at ASC "
		}
		if payload.CursorSortValue != nil {
			cursorTime, err := time.Parse(time.RFC3339Nano, *payload.CursorSortValue)
			if err != nil {
				return nil, fmt.Errorf("failed to parse cursor sort value: %w", err)
			}

			args["cursor_sort_value"] = cursorTime
			if *payload.Order == model.OrderDesc {
				stmt += `
					AND (
					    p.created_at <= @cursor_sort_value
					)`
			} else {
				stmt += `
					AND (
					    p.created_at >= @cursor_sort_value
					)`
			}
		}
	} else {
		if payload.Order == nil || *payload.Order == model.OrderDesc {
			orderStmt += " ORDER BY p.upvotes DESC, p.created_at DESC "
		} else {
			orderStmt += " ORDER BY p.upvotes ASC, p.created_at DESC "
		}
		if payload.CursorSortValue != nil {
			upvotes, err := strconv.Atoi(*payload.CursorSortValue)
			if err != nil {
				return nil, fmt.Errorf("Failed to convert sort value to int: %w", err)
			}
			args["cursor_sort_value"] = upvotes
			args["cursor_created_at"] = payload.CursorCreatedAt
			if *payload.Order == model.OrderDesc {
				stmt += ` AND (
					p.upvotes < @cursor_sort_value
					OR (
						p.upvotes = @cursor_sort_value
						AND p.created_at <= @cursor_created_at
					)
				)`
			} else {
				stmt += ` AND (
					p.upvotes > @cursor_sort_value
					OR (
						p.upvotes = @cursor_sort_value
						AND p.created_at >= @cursor_created_at
					)
				)`
			}
		}
	}

	stmt += " GROUP BY p.id"
	stmt += orderStmt
	stmt += " LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to query get posts for user id %s: %w", userID, err)
	}

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[post.PopulatedPost])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.CursorPaginatedResponse[post.PopulatedPost]{
				Data:            []post.PopulatedPost{},
				CursorSortValue: *payload.CursorSortValue,
				CursorCreatedAt: *payload.CursorCreatedAt,
				HasMore:         false,
			}, nil
		}

		return nil, fmt.Errorf("failed to collect rows from tables: todos for user_id=%s : %w", userID, err)
	}

	if len(posts) < limit+1 {
		length := len(posts)
		return &model.CursorPaginatedResponse[post.PopulatedPost]{
			Data:            posts[:length],
			CursorSortValue: *payload.CursorSortValue,
			CursorCreatedAt: *payload.CursorCreatedAt,
			HasMore:         false,
		}, nil
	}
	var nextCursorSortValue string
	var nextCursorCreatedAt time.Time
	if payload.Sort == nil || *payload.Sort == model.SortByCreatedAt {
		nextCursorSortValue = posts[limit].CreatedAt.Format(time.RFC3339Nano)
		nextCursorCreatedAt = posts[limit].CreatedAt
	} else {
		nextCursorSortValue = strconv.Itoa(posts[limit].Upvotes)
		nextCursorCreatedAt = posts[limit].CreatedAt
	}

	HasMore := true

	return &model.CursorPaginatedResponse[post.PopulatedPost]{
		Data:            posts[:limit],
		CursorSortValue: nextCursorSortValue,
		CursorCreatedAt: nextCursorCreatedAt,
		HasMore:         HasMore,
	}, nil

}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *UserRepository) GetUsers(ctx context.Context, payload *user.GetUsersQuery) (*model.CursorPaginatedResponse[user.MiniUser], error) {

	stmt := `
		SELECT u.id, u.username, u.display_name, u.avatar_key, u.follower_count, u.bio, u.created_at
		FROM users u
	`
	conditions := []string{}
	args := pgx.NamedArgs{}

	if payload.Name != nil {
		args["name"] = *payload.Name + "%"
		conditions = append(conditions, "(username ILIKE @name OR display_name ILIKE @name)")
	}

	if payload.DateJoinedStart != nil {
		args["date_joined_start"] = *payload.DateJoinedStart
		conditions = append(conditions, "u.created_at >= @date_joined_start")
	}

	if payload.DateJoinedEnd != nil {
		args["date_joined_end"] = *payload.DateJoinedEnd
		conditions = append(conditions, "u.created_at <= @date_joined_end")
	}

	orderStmt := ""

	if payload.Sort == nil || *payload.Sort == model.SortByFollowerCount {
		if payload.Order == nil || *payload.Order == model.OrderDesc {
			orderStmt += " ORDER BY u.follower_count DESC, u.created_at DESC"
		} else {
			orderStmt += " ORDER BY u.follower_count ASC, u.created_at ASC"
		}

		if payload.CursorSortValue != nil {
			cursorFollowerCount, err := strconv.Atoi(*payload.CursorSortValue)
			if err != nil {
				return nil, fmt.Errorf("Failed to convert to int: %w", err)
			}
			args["cursor_follower_count"] = cursorFollowerCount
			args["cursor_created_at"] = payload.CursorCreatedAt
			if payload.Order == nil || *payload.Order == model.OrderDesc {
				conditions = append(conditions, "((u.follower_count < @cursor_follower_count) OR (u.follower_count = @cursor_follower_count AND u.created_at <= @cursor_created_at))")

			} else {
				conditions = append(conditions, "((u.follower_count > @cursor_follower_count) OR (u.follower_count = @cursor_follower_count AND u.created_at >= @cursor_created_at))")
			}
		}
	} else {
		if payload.Order == nil || *payload.Order == model.OrderDesc {
			orderStmt += " ORDER BY u.created_at DESC"
		} else {
			orderStmt += " ORDER BY u.created_at ASC"
		}

		if payload.CursorSortValue != nil {
			args["cursor_created_at"] = payload.CursorCreatedAt
			if payload.Order == nil || *payload.Order == model.OrderDesc {
				conditions = append(conditions, "u.created_at <= @cursor_created_at")
			} else {
				conditions = append(conditions, "u.created_at >= @cursor_created_at")
			}
		}
	}
	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	stmt += " " + orderStmt

	limit := 20
	args["limit_plus_one"] = limit + 1

	stmt += " LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to process get users query: %w", err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[user.MiniUser])

	if err != nil {
		return nil, fmt.Errorf("Faild to collect rows: %w", err)
	}

	if len(users) < limit+1 {
		return &model.CursorPaginatedResponse[user.MiniUser]{
			Data:            users,
			CursorSortValue: *payload.CursorSortValue,
			CursorCreatedAt: *payload.CursorCreatedAt,
			HasMore:         false,
		}, nil
	}

	if payload.Sort == nil || *payload.Sort == model.SortByFollowerCount {
		nextCursorCreatedAt := users[limit].CreatedAt
		nextFollowerCount := strconv.Itoa(users[limit].FollowerCount)
		return &model.CursorPaginatedResponse[user.MiniUser]{
			Data:            users[:limit],
			CursorSortValue: nextFollowerCount,
			CursorCreatedAt: nextCursorCreatedAt,
			HasMore:         true,
		}, nil
	}

	nextCursorCreatedAt := users[limit].CreatedAt

	return &model.CursorPaginatedResponse[user.MiniUser]{
		Data:            users[:limit],
		CursorSortValue: nextCursorCreatedAt.Format(time.RFC3339Nano),
		CursorCreatedAt: nextCursorCreatedAt,
		HasMore:         true,
	}, nil

}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
