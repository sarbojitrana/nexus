package community

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//-------------------------------------------------------------------------------------------

type CreateCommunityPayload struct {
	Name        string  `json:"name" validate:"required,max=50"`
	Slug        string  `json:"slug" validate:"required,max=50"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	AvatarKey   *string `json:"avatarKey"`
	BannerKey   *string `json:"bannerKey"`
}

func (p *CreateCommunityPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.Description == nil {
		defaultDescription := ""
		p.Description = &defaultDescription
	}

	return nil
}

//-------------------------------------------------------------------------------------------

type UpdateCommunityPayload struct {
	ID          uuid.UUID  `json:"id" validate:"required,uuid"`
	AdminID     *uuid.UUID `json:"adminId" validate:"required,uuid"`
	Name        *string    `json:"name" validate:"omitempty,max=50"`
	Slug        *string    `json:"slug" validate:"omitempty,max=50"`
	Description *string    `json:"description" validate:"omitempty,max=1000"`
	AvatarKey   *string    `json:"avatarKey"`
	BannerKey   *string    `json:"bannerKey"`
}

func (p *UpdateCommunityPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}

//-------------------------------------------------------------------------------------------

type DeleteCommunityPayload struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (p *DeleteCommunityPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}

//-------------------------------------------------------------------------------------------

type CommunityFollowPayload struct {
	CommunityID uuid.UUID `json:"communityId" validate:"required,uuid"`
	FollowerID  uuid.UUID `json:"followerId" validate:"required,uuid"`
}

func (p *CommunityFollowPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}

//-------------------------------------------------------------------------------------------

type DeleteCommunityFollowPayload struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (p *DeleteCommunityFollowPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}

//-------------------------------------------------------------------------------------------

type GetCommunityByIDPayload struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (p *GetCommunityByIDPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}

//-------------------------------------------------------------------------------------------

type GetCommunitiesPayload struct {
	Name             *string `json:"name" validate:"omitempty,max=50"`
	NextCursor       *string `json:"nextCursor"`
	Sort             *string `json:"sort" validate:"omitempty,oneof=created_at members_count posts_count"`
	Order            *string `json:"order" validate:"omitempty,oneof=asc desc"`
	DateCreatedStart *string `json:"dateCreatedStart"`
	DateCreatedEnd   *string `json:"dateCreatedEnd"`
}

func (p *GetCommunitiesPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.Name == nil {
		defaultName := ""
		p.Name = &defaultName
	}

	if p.NextCursor == nil {
		defaultCursor := ""
		p.NextCursor = &defaultCursor
	}
	if p.Sort == nil {
		defaultSort := "members_count"
		p.Sort = &defaultSort
	}

	if p.Order == nil {
		defaultOrder := "desc"
		p.Order = &defaultOrder
	}
	return nil
}
