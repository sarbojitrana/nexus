package user

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//-------------------------------------------------------------------------------------------

type CreateUserPayload struct {
	Username    string  `json:"username" validate:"required"`
	DisplayName string  `json:"displayName" validate:"required"`
	EmailID     string  `json:"emailId" validate:"required"`
	ClerkID     string  `json:"clerkId" validate:"required"`
	Bio         *string `json:"bio" validate:"omitempty,max=1000"`
	AvatarKey   *string `json:"avatarKey"`
	BannerKey   *string `json:"bannerKey"`
}

func (p *CreateUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type UpdateUserPayload struct {
	Username    *string `json:"username" validate:"omitempty,max=50"`
	EmailID     *string `json:"emailId" validate:"omitempty"`
	DisplayName *string `json:"displayName" validate:"omitempty,max=50"`
	Bio         *string `json:"bio" validate:"omitempty,max=1000"`
	AvatarKey   *string `json:"avatarKey" validate:"omitempty"`
	BannerKey   *string `json:"bannerKey" validate:"omitempty"`
}

func (p *UpdateUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type GetUserByIDPayload struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (p *GetUserByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type GetUsersPayload struct {
	NextCursor      *string `json:"nextCursor" validate:"omitempty"`
	SortBy          *string `json:"sortBy" validate:"omitempty,oneof=created_at display_name follower_count following_count posts_count"`
	Order           *string `json:"order" validate:"omitempty,oneof=asc desc"`
	Name            *string `json:"name" validate:"omitempty"`
	DateJoinedStart *string `json:"dateJoinedStart" validate:"omitempty"`
	DateJoinedEnd   *string `json:"dateJoinedEnd" validate:"omitempty"`
}

func (p *GetUsersPayload) Validate() error {
	validate := validator.New()

	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.NextCursor == nil {
		defaultCursor := ""
		p.NextCursor = &defaultCursor
	}

	if p.SortBy == nil {
		defaultSort := "follower_count"
		p.SortBy = &defaultSort
	}

	if p.Order == nil {
		defaultOrder := "desc"
		p.Order = &defaultOrder
	}

	if p.Name == nil {
		defaultName := ""
		p.Name = &defaultName
	}

	return nil
}

//-------------------------------------------------------------------------------------------

type GetPostsByUserIDPayload struct {
	NextCursor       *string    `json:"nextCursor"`
	SortBy           *string    `json:"sortBy" validate:"omitempty,oneof=created_at upvotes"`
	Order            *string    `json:"order" validate:"omitempty,oneof=asc desc"`
	DateCreatedStart *time.Time `json:"dateCreatedStart"`
	DateCreatedEnd   *time.Time `json:"dateCreatedEnd"`
}

func (p *GetPostsByUserIDPayload) Validate() error {
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
		defaultOrder := "desc"
		p.Order = &defaultOrder
	}

	return nil
}
