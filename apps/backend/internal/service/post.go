package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/post"
	"github.com/sarbojitrana/nexus/internal/repository"
	"github.com/sarbojitrana/nexus/internal/server"
)

type PostService struct {
	server *server.Server
	repo   *repository.PostRepository
}

func NewPostService(s *server.Server, repo *repository.PostRepository) *PostService {
	return &PostService{
		server: s,
		repo:   repo,
	}
}

func (s *PostService) CreatePost(ctx context.Context, userID string, payload *post.CreatePostPayload) (*post.Post, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	created, err := s.repo.CreatePost(ctx, userID, payload)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to create post")
		return nil, err
	}

	logger.Info().Str("event", "post_created").Str("post_id", created.ID.String()).Str("user_id", userID).Msg("post created")
	return created, nil
}

func (s *PostService) UpdatePost(ctx context.Context, userID string, postID uuid.UUID, payload *post.UpdatePostByIDPayload) (*post.Post, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	updated, err := s.repo.UpdatePostByID(ctx, userID, postID, payload)
	if err != nil {
		logger.Error().Err(err).Str("post_id", postID.String()).Str("user_id", userID).Msg("failed to update post")
		return nil, err
	}

	logger.Info().Str("event", "post_updated").Str("post_id", postID.String()).Str("user_id", userID).Msg("post updated")
	return updated, nil
}

func (s *PostService) DeletePost(ctx context.Context, userID string, postID uuid.UUID) error {
	logger := middleware.GetLoggerFromContext(ctx)

	if err := s.repo.DeletePostByID(ctx, userID, &post.DeletePostByIDPayload{ID: postID}); err != nil {
		logger.Error().Err(err).Str("post_id", postID.String()).Str("user_id", userID).Msg("failed to delete post")
		return err
	}

	logger.Info().Str("event", "post_deleted").Str("post_id", postID.String()).Str("user_id", userID).Msg("post deleted")
	return nil
}

func (s *PostService) GetPostByID(ctx context.Context, viewerID *string, postID uuid.UUID) (*post.PopulatedPost, error) {
	p, err := s.repo.GetPostByID(ctx, viewerID, postID)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("post_id", postID.String()).Msg("failed to fetch post")
		return nil, err
	}
	return p, nil
}

func (s *PostService) GetComments(ctx context.Context, viewerID *string, postID uuid.UUID, query *post.GetCommentsByPostIDQuery) (*model.CursorPaginatedResponse[post.PopulatedPost], error) {
	res, err := s.repo.GetCommentsByPostID(ctx, viewerID, postID, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("post_id", postID.String()).Msg("failed to fetch comments")
		return nil, err
	}
	return res, nil
}

func (s *PostService) GetReplies(ctx context.Context, viewerID *string, commentID uuid.UUID, query *post.GetRepliesByCommentIDQuery) (*model.OffsetPaginatedResponse[post.PopulatedPost], error) {
	res, err := s.repo.GetRepliesByCommentID(ctx, viewerID, commentID, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("comment_id", commentID.String()).Msg("failed to fetch replies")
		return nil, err
	}
	return res, nil
}

func (s *PostService) React(ctx context.Context, userID string, postID uuid.UUID, payload *post.ReactToPostPayload) error {
	logger := middleware.GetLoggerFromContext(ctx)

	if err := s.repo.ReactToPost(ctx, userID, postID, payload); err != nil {
		logger.Error().Err(err).Str("post_id", postID.String()).Str("user_id", userID).Msg("failed to react to post")
		return err
	}

	logger.Info().Str("event", "post_reacted").Str("post_id", postID.String()).Str("user_id", userID).Str("reaction", string(payload.Reaction)).Msg("post reaction recorded")
	return nil
}

func (s *PostService) GetFeed(ctx context.Context, viewerID *string, query *post.GetPostsQuery) (*post.GetPostsQueryResponse, error) {
	res, err := s.repo.GetPosts(ctx, viewerID, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Msg("failed to fetch feed")
		return nil, err
	}
	return res, nil
}
