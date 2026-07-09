package user

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
)

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetUserByIDQuery struct {
	ID uuid.UUID `query:"id" validate:"required,uuid"`
}

func (p *GetUserByIDQuery) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetUsersQuery struct {
	CursorSortValue *string      `json:"cursorSortValue"`
	CursorCreatedAt *time.Time   `json:"cursorCreatedAt"`
	Sort            *model.Sort  `json:"sort" validate:"omitempty,oneof=created_at follower_count"`
	Order           *model.Order `json:"order" validate:"omitempty,oneof=asc desc"`
	Name            *string      `json:"name" validate:"omitempty"`
	DateJoinedStart *time.Time   `json:"dateJoinedStart" validate:"omitempty"`
	DateJoinedEnd   *time.Time   `json:"dateJoinedEnd" validate:"omitempty"`
}

func (p *GetUsersQuery) Validate() error {
	validate := validator.New()

	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.Sort == nil {
		defaultSortBy := model.SortByFollowerCount
		p.Sort = &defaultSortBy
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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetPostsByUserIDPayload struct {
	CursorSortValue  *string      `query:"cursorSortValue"`
	CursorCreatedAt  *time.Time   `query:"cursorCreatedAt"`
	Sort             *model.Sort  `query:"sort" validate:"omitempty,oneof=created_at upvotes"`
	Order            *model.Order `query:"order" validate:"omitempty,oneof=asc desc"`
	DateCreatedStart *time.Time   `query:"dateCreatedStart"`
	DateCreatedEnd   *time.Time   `query:"dateCreatedEnd"`
}

func (p *GetPostsByUserIDPayload) Validate() error {
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
