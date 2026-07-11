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
	AuthorID                   *uuid.UUID   `query:"authorId" validate:"omitempty,uuid"`
	CommunityID                *uuid.UUID   `query:"communityId" validate:"omitempty,uuid"`
	ParentPostID               *uuid.UUID   `query:"parentPostId" validate:"omitempty,uuid"`
	NextTrendingPostValue      int64        `query:"nextTrendingPostValue"`
	NextTrendingPostCreatedAt  time.Time    `query:"nextTrendingPostCreatedAt"`
	NextFollowingPostValue     int64        `query:"nextFollowingPostValue"`
	NextFollowingPostCreatedAt time.Time    `query:"nextFollowingPostCreatedAt"`
	Sort                       *model.Sort  `query:"sort" validate:"omitempty,oneof=created_at upvotes"`
	Order                      *model.Order `query:"order" validate:"omitempty,oneof=asc desc"`
	DateCreatedStart           *time.Time   `query:"dateCreatedStart"`
	DateCreatedEnd             *time.Time   `query:"dateCreatedEnd"`
}

func (p *GetPostsQuery) Validate() error {
	validate := validator.New()

	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.Sort == nil {
		defaultSortBy := model.SortByCreatedAt
		p.Sort = &defaultSortBy

	}

	if p.Order == nil {
		defaultOrder := model.OrderDesc
		p.Order = &defaultOrder
	}
	return nil
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
	VoteType VoteType  `json:"voteType" validate:"required,oneof=-1 1"`
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
	Reaction VoteType `json:"reaction" validate:"required, oneof= upvote downvote"`
}

func (p *ReactToPostPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
