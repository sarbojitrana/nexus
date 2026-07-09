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
	"github.com/sarbojitrana/nexus/internal/server"
)

type PostRepository struct {
	server *server.Server
}

func NewPostRepository(server *server.Server) *PostRepository {
	return &PostRepository{
		server: server,
	}
}

func (r *PostRepository) CreatePost(ctx context.Context, userID uuid.UUID, payload *post.CreatePostPayload) (*post.Post, error) {
	stmt := `
		INSERT INTO posts (
			author_id,
			community_id,
			parent_post_id,
			post_type,
			title,
			content
		)
		VALUES(
			@user_id,
			@community_id,
			@parent_post_id,
			@post_type,
			@title,
			@content
		)
		RETURNING
		*
	`
	args := pgx.NamedArgs{
		"user_id":        userID,
		"community_id":   payload.CommunityID,
		"parent_post_id": payload.ParentPostID,
		"post_type":      payload.PostType,
		"title":          payload.Title,
		"content":        payload.Content,
	}

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to create post for user_id %s: %w", userID, err)
	}
	post, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[post.Post])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse the post to a struct for user_id %s: %w", userID, err)
	}

	return &post, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *PostRepository) DeletePostByID(ctx context.Context, userID uuid.UUID, payload *post.DeletePostByIDPayload) error {
	stmt := `
		DELETE FROM posts
		WHERE author_id = @author_id
		AND id = @id
	`
	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"author_id": userID,
		"id":        payload.ID,
	})

	if err != nil {
		return fmt.Errorf("Failed to delete post for user_id %s: %w", userID, err)
	}

	if result.RowsAffected() == 0 {
		code := "POST NOT FOUND"
		return errs.NewNotFoundError("post not found", false, &code)
	}
	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *PostRepository) UpdatePostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID, payload *post.UpdatePostByIDPayload) (*post.Post, error) {
	stmt := `
		UPDATE posts
	`
	args := pgx.NamedArgs{
		"post_id": postID,
		"user_id": userID,
	}
	setClauses := []string{}

	if payload.Content != nil {
		args["content"] = payload.Content
		setClauses = append(setClauses, "content = @content")
	}

	if payload.Title != nil {
		args["title"] = payload.Title
		setClauses = append(setClauses, "title = @title")
	}

	if len(setClauses) == 0 {
		code := "NOTHING TO UPDATE"
		return nil, errs.NewBadRequestError("No fields sent to update", false, &code, nil, nil)
	}

	stmt += " SET " + strings.Join(setClauses, ", ") + " WHERE id = @post_id AND author_id = @user_id RETURNING *"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to update post of post_id %s and user_id %s: %w", postID, userID, err)
	}

	post, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[post.Post])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse the row to struct for post_id %s and user_id %s: %w", postID, userID, err)
	}

	return &post, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *PostRepository) GetPostByID(ctx context.Context, postID uuid.UUID) (*post.PopulatedPost, error) {
	stmt := `
		SELECT p.*,
		COALESCE(
			json_agg(
				to_jsonb(camel(pm))
				ORDER BY
					pm.created_at DESC,
					pm.id DESC
			) FILTER(
				WHERE pm.id IS NOT NULL 
			),
			'[]' :: JSONB
		) AS post_media

		FROM posts p
		LEFT JOIN post_media pm ON pm.post_id = p.id
		WHERE p.id = @post_id
		GROUP BY p.id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"post_id": postID,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to get post of post_id %s: %w", postID, err)
	}

	post, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[post.PopulatedPost])

	if err != nil {
		return nil, fmt.Errorf("Failed to convert rows to struct: %w", err)
	}
	return &post, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *PostRepository) GetCommentsByPostID(ctx context.Context, postID uuid.UUID, payload *post.GetCommentsByPostIDQuery) (*model.CursorPaginatedResponse[post.PopulatedPost], error) {
	stmt := `
		SELECT c.*,
		COALESCE(
			json_agg(
				to_jsonb(camel(pm))
				ORDER BY
					pm.created_at DESC,
					pm.id DESC
			) FILTER(
				WHERE pm.id IS NOT NULL 
			),
			'[]' :: JSONB
		) AS post_media

		FROM posts c
		LEFT JOIN post_media pm ON pm.post_id = c.id
	`
	limit := 20
	args := pgx.NamedArgs{
		"post_id":        postID,
		"limit_plus_one": limit + 1,
	}

	conditions := []string{}

	if payload.DateCreatedStart != nil {
		conditions = append(conditions, "c.created_at >= @date_created_start")
		args["date_created_start"] = *payload.DateCreatedStart
	}

	if payload.DateCreatedEnd != nil {
		conditions = append(conditions, "c.created_at <= @date_created_end")
		args["date_created_end"] = *payload.DateCreatedEnd
	}

	stmt += " WHERE c.parent_post_id = @post_id"
	if len(conditions) != 0 {
		stmt += " AND " + strings.Join(conditions, " AND ")
	}

	orderStmt := ""

	if payload.Sort == nil || *payload.Sort == model.SortByPopularity {
		if payload.Order == nil || *payload.Order == model.OrderDesc {
			orderStmt += " ORDER BY c.upvotes + c.downvotes + c.comment_count DESC, c.created_at DESC "
		} else {
			orderStmt += " ORDER BY c.upvotes + c.downvotes + c.comment_count ASC, c.created_at ASC "
		}

		if payload.CursorSortValue != nil {
			cursorSortValue, err := strconv.Atoi(*payload.CursorSortValue)
			if err != nil {
				return nil, fmt.Errorf("Failed to convert sort parameter to int: %w", err)
			}
			args["cursor_sort_value"] = cursorSortValue
			args["cursor_created_at"] = payload.CursorCreatedAt

			if *payload.Order == model.OrderDesc {
				stmt += ` AND (
					c.upvotes + c.comment_count + c.downvotes < @cursor_sort_value
					OR (
						c.upvotes + c.comment_count + c.downvotes = @cursor_sort_value
						AND c.created_at <= @cursor_created_at
					)
				)`
			} else {
				stmt += ` AND (
					c.upvotes + c.comment_count + c.downvotes > @cursor_sort_value
					OR (
						c.upvotes + c.comment_count + c.downvotes = @cursor_sort_value
						AND c.created_at >= @cursor_created_at
					)
				)`
			}
		}

	} else {
		if payload.Order == nil || *payload.Order == model.OrderDesc {
			orderStmt += " ORDER BY c.created_at DESC "
		} else {
			orderStmt += " ORDER BY c.created_at ASC "
		}

		if payload.CursorSortValue != nil {
			cursorTime, err := time.Parse(time.RFC3339Nano, *payload.CursorSortValue)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse cursor sort value: %w", err)
			}
			args["cursor_sort_value"] = cursorTime
			if *payload.Order == model.OrderDesc {
				stmt += `
					AND (
						c.created_at <= @cursor_sort_value
					)
				`
			} else {
				stmt += `
					AND (
						c.created_at >= @cursor_sort_value
					)
				`
			}
		}
	}

	stmt += " GROUP BY c.id"
	stmt += orderStmt
	stmt += " LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to query to get comments for post_id %s: %w", postID, err)
	}

	comments, err := pgx.CollectRows(rows, pgx.RowToStructByName[post.PopulatedPost])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.CursorPaginatedResponse[post.PopulatedPost]{
				Data:    []post.PopulatedPost{},
				CursorSortValue: *payload.CursorSortValue,
				CursorCreatedAt: *payload.CursorCreatedAt,
				HasMore: false,
			}, nil
		}

		return nil, fmt.Errorf("Failed to collect rows from tables for post_id=%s: %w", postID, err)

	}

	if len(comments) < limit+1 {
		length := len(comments)
		return &model.CursorPaginatedResponse[post.PopulatedPost]{
			Data:    comments[:length],
			CursorSortValue: *payload.CursorSortValue,
			CursorCreatedAt: *payload.CursorCreatedAt,
			HasMore: false,
		}, nil
	}

	var nextCursorSortValue string
	var nextCursorCreatedAt time.Time

	if payload.Sort == nil || *payload.Sort == model.SortByPopularity {
		nextCursorSortValue = strconv.Itoa((comments[limit].Upvotes + comments[limit].CommentCount + comments[limit].Downvotes))
		nextCursorCreatedAt = comments[limit].CreatedAt
	} else {
		nextCursorSortValue = comments[limit].CreatedAt.Format(time.RFC3339Nano)
		nextCursorCreatedAt = comments[limit].CreatedAt
	}

	HasMore := true

	return &model.CursorPaginatedResponse[post.PopulatedPost]{
		Data:    comments[:limit],
		CursorSortValue: nextCursorSortValue,
		CursorCreatedAt: nextCursorCreatedAt,
		HasMore: HasMore,
	}, nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (r *PostRepository) GetRepliesByCommentID(ctx context.Context, commentID uuid.UUID, query *post.GetRepliesByCommentIDQuery) (*model.OffsetPaginatedResponse[post.PopulatedPost], error) {

	stmt := `
		SELECT r.*,
		COALESCE(
			json_agg(
				to_jsonb(camel(pm))
				ORDER BY pm.created_at DESC
			) FILTER(
				WHERE pm.id IS NOT NULL
			),
			'[]' :: jsonb
		) AS post_media
		FROM posts r
		LEFT JOIN post_media pm
		ON pm.post_id = r.id
		WHERE r.parent_post_id = @comment_id 
		GROUP BY r.id
		ORDER BY r.upvotes + r.downvotes + r.comment_count DESC, r.created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"comment_id": commentID,
		"limit":      *query.Limit,
		"offset":     (*query.Page - 1) * (*query.Limit),
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to query to get replies for comment_id %s: %w", commentID, err)
	}

	replies, err := pgx.CollectRows(rows, pgx.RowToStructByName[post.PopulatedPost])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.OffsetPaginatedResponse[post.PopulatedPost]{
				Data:       []post.PopulatedPost{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("Failed to collect rows for comment_id %s: %w", commentID, err)
	}

	countStmt := `
		SELECT COUNT(*)
		FROM posts r
		WHERE r.parent_post_id = @comment_id
	`

	var total int

	err = r.server.DB.Pool.QueryRow(ctx, countStmt, pgx.NamedArgs{
		"comment_id": commentID,
	}).Scan(&total)

	if err != nil {
		return nil, fmt.Errorf("Failed to count total replies for comment_id %s: %w", commentID, err)
	}

	return &model.OffsetPaginatedResponse[post.PopulatedPost]{
		Data:       replies,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil

}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
