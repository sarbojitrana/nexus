package post

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
)

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type CreatePostPayload struct {
	CommunityID  *uuid.UUID `json:"communityId" validate:"omitempty,uuid"`
	ParentPostID *uuid.UUID `json:"parentPostId" validate:"omitempty,uuid"`
	PostType     *PostType  `json:"postType" validate:"oneof=comment post"`
	Title        *string    `json:"title"`
	Content      *string    `json:"content"`
}

func (p *CreatePostPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type DeletePostByIDPayload struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (p *DeletePostByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type UpdatePostByIDPayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func (p *UpdatePostByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetPostsQuery struct {
	ReferenceTime             *time.Time `query:"referenceTime"`
	TrendingLimit             *int       `query:"trendingLimit"`
	FollowingUsersLimit       *int       `query:"followingUsersLimit"`
	FollowingCommunitiesLimit *int       `query:"followingCommunitiesLimit"`

	TrendingCursorValue                 *float64   `query:"trendingCursorValue"`
	TrendingCursorCreatedAt             *time.Time `query:"trendingCursorCreatedAt"`

	FollowingUsersCursorCreatedAt       *time.Time `query:"followingUsersCursorCreatedAt"`
	FollowingCommunitiesCursorCreatedAt *time.Time `query:"followingCommunitiesCursorCreatedAt"`
}

func (p *GetPostsQuery) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetPostsQueryResponse struct {
	ReferenceTime                           time.Time       `json:"referenceTime"`
	
	TrendingPosts                           []PopulatedPost `json:"trendingPosts"`
	NextTrendingCursorValue                 *float64        `json:"nextTrendingCursorValue"`
	NextTrendingCursorCreatedAt             *time.Time      `json:"nextTrendingCursorCreatedAt"`
	HasMoreTrending                         bool            `json:"hasMoreTrending"`

	FollowingUsersPosts                     []PopulatedPost `json:"followingUsersPosts"`
	NextFollowingUsersCursorCreatedAt       *time.Time      `json:"nextFollowingUsersCursorCreatedAt"`
	HasMoreFollowingUsers                   bool            `json:"hasMoreFollowingUsers"`

	FollowingCommunitiesPosts               []PopulatedPost `json:"followingCommunitiesPosts"`
	NextFollowingCommunitiesCursorCreatedAt *time.Time      `json:"nextFollowingCommunitiesCursorCreatedAt"`
	HasMoreFollowingCommunities             bool            `json:"hasMoreFollowingCommunities"`
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type CreatePostMediaPayload struct {
	PostID     uuid.UUID `json:"postId" validate:"required,uuid"`
	StorageKey string    `json:"storageKey" validate:"required"`
	FileSize   int64     `json:"fileSize" validate:"required"`
	MimeType   string    `json:"mimeType" validate:"required"`
}

func (p *CreatePostMediaPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type DeletePostMediaPayload struct {
	ID     uuid.UUID `json:"id" validate:"required,uuid"`
	UserID uuid.UUID `json:"userId" validate:"required,uuid"`
}

func (p *DeletePostMediaPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetPostMediaQuery struct {
	PostID uuid.UUID `query:"postId" validate:"required,uuid"`
}

func (p *GetPostMediaQuery) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type CreatePostVotePayload struct {
	PostID   uuid.UUID `json:"postId" validate:"required,uuid"`
	VoteType VoteType  `json:"voteType" validate:"required,oneof=upvote downvote"`
}

func (p *CreatePostVotePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type DeletePostVotePayload struct {
	ID     uuid.UUID `json:"id" validate:"required,uuid"`
	UserID uuid.UUID `json:"userId" validate:"required,uuid"`
}

func (p *DeletePostVotePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetCommentsByPostIDQuery struct {
	CursorSortValue  *string      `query:"cursorSortValue"`
	CursorCreatedAt  *time.Time   `query:"cursorCreatedAt"`
	Sort             *model.Sort  `query:"sort" validate:"omitempty,oneof=created_at upvotes"`
	Order            *model.Order `query:"order" validate:"omitempty,oneof=asc desc"`
	DateCreatedStart *time.Time   `query:"dateCreatedStart" validate:"omitempty"`
	DateCreatedEnd   *time.Time   `query:"dateCreatedEnd" validate:"omitempty"`
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetRepliesByCommentIDQuery struct {
	Page  *int `query:"page" validate:"omitempty,min=1"`
	Limit *int `query:"limit" validate:"omitempty,min=1,max=100"`
}

func (q *GetRepliesByCommentIDQuery) Validate() error {

	validate := validator.New()

	if err := validate.Struct(q); err != nil {
		return err
	}

	if q.Page == nil {
		defaultPage := 1
		q.Page = &defaultPage
	}

	if q.Limit == nil {
		defaultLimit := 5
		q.Limit = &defaultLimit
	}

	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type ReactToPostPayload struct {
	Reaction VoteType `json:"reaction" validate:"required,oneof=upvote downvote"`
}

func (p *ReactToPostPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
