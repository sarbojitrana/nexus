package user

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
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
	NextCursor      *model.Cursor `json:"nextCursor" validate:"omitempty"`
	SortBy          *model.SortBy `json:"sortBy" validate:"omitempty,oneof=created_at follower_count"`
	Order           *model.Order  `json:"order" validate:"omitempty,oneof=asc desc"`
	Name            *string       `json:"name" validate:"omitempty"`
	DateJoinedStart *time.Time    `json:"dateJoinedStart" validate:"omitempty"`
	DateJoinedEnd   *time.Time    `json:"dateJoinedEnd" validate:"omitempty"`
}

func (p *GetUsersPayload) Validate() error {
	validate := validator.New()

	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.SortBy == nil {
		defaultSortBy := model.SortByFollowerCount
		p.SortBy = &defaultSortBy
	}

	if p.Order == nil {
		defaultOrder := model.OrderDesc
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
	NextCursor       *model.Cursor `json:"nextCursor" valdate:"omitempty"`
	SortBy           *model.SortBy `json:"sortBy" validate:"omitempty,oneof=created_at upvotes"`
	Order            *model.Order  `json:"order" validate:"omitempty,oneof=asc desc"`
	DateCreatedStart *time.Time    `json:"dateCreatedStart"`
	DateCreatedEnd   *time.Time    `json:"dateCreatedEnd"`
}

func (p *GetPostsByUserIDPayload) Validate() error {
	validate := validator.New()

	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.SortBy == nil {
		defaultSortBy := model.SortByCreatedAt
		p.SortBy = &defaultSortBy
	}

	if p.Order == nil {
		defaultOrder := model.OrderDesc
		p.Order = &defaultOrder

	}

	return nil
}
