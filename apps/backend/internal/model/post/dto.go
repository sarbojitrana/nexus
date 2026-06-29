package post

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//-------------------------------------------------------------------------------------------

type CreatePostPayload struct {
	AuthorID     uuid.UUID  `json:"authorId" validate:"required,uuid"`
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

//-------------------------------------------------------------------------------------------

type DeletePostByIDPayload struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (p *DeletePostByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type UpdatePostByIDPayload struct {
	ID      uuid.UUID `json:"id" validate:"required,uuid"`
	Title   *string   `json:"title"`
	Content *string   `json:"content"`
}

func (p *UpdatePostByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type GetPostByIDPayload struct {
	ID uuid.UUID
}

//-------------------------------------------------------------------------------------------

type GetPostsPayload struct {
	AuthorID         *uuid.UUID `json:"authorId" validate:"omitempty,uuid"`
	CommunityID      *uuid.UUID `json:"communityId" validate:"omitempty,uuid"`
	ParentPostID     *uuid.UUID `json:"parentPostId" validate:"omitempty,uuid"`
	NextCursor       *string    `json:"nextCursor"`
	SortBy           *string    `json:"sortBy" validate:"omitempty,oneof=created_at upvotes downvotes"`
	Order            *string    `json:"order" validate:"omitempty,oneof=asc desc"`
	DateCreatedStart *string    `json:"dateCreatedStart"`
	DateCreatedEnd   *string    `json:"dateCreatedEnd"`
}

func (p *GetPostsPayload) Validate() error {
	validate := validator.New()

	if err := validate.Struct(p); err != nil {
		return err
	}
	if p.NextCursor == nil {
		defaultCursor := ""
		p.NextCursor = &defaultCursor
	}

	if p.SortBy == nil {
		defaultSort := "created_at"
		p.SortBy = &defaultSort
	}

	if p.Order == nil {
		defaultOrder := "asc"
		p.Order = &defaultOrder
	}

	return nil
}

//-------------------------------------------------------------------------------------------

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

//-------------------------------------------------------------------------------------------

type DeletePostMediaPayload struct {
	ID     uuid.UUID `json:"id" validate:"required,uuid"`
	UserID uuid.UUID `json:"userId" validate:"required,uuid"`
}

func (p *DeletePostMediaPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type GetPostMediaPayload struct {
	PostID uuid.UUID `json:"postId" validate:"required,uuid"`
}

func (p *GetPostMediaPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type CreatePostVotePayload struct {
	PostID   uuid.UUID `json:"postId" validate:"required,uuid"`
	VoteType VoteType  `json:"voteType" validate:"required,oneof=-1 1"`
}

func (p *CreatePostVotePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type DeletePostVotePayload struct {
	ID     uuid.UUID `json:"id" validate:"required,uuid"`
	UserID uuid.UUID `json:"userId" validate:"required,uuid"`
}

func (p *DeletePostVotePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
