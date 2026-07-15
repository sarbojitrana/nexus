package follow

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
)

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type FollowCommunityPayload struct {
	CommunityID uuid.UUID `json:"communityId" validate:"required,uuid"`
}

func (p *FollowCommunityPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type UnFollowCommunityPayload struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (p *UnFollowCommunityPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type FollowUserPayload struct {
	FollowingID uuid.UUID `json:"followingId" validate:"required,uuid"`
}

func (p *FollowUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type UnFollowUserPayload struct {
	FollowingID uuid.UUID `json:"followingId" validate:"required,uuid"`
}

func (p *UnFollowUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetFollowersQuery struct {
	CursorCreatedAt *time.Time   `query:"cursorCreatedAt"`
	Order           *model.Order `query:"order" validate:"omitempty,oneof=asc desc"`
}

func (p *GetFollowersQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	if p.Order == nil {
		defaultOrder := model.OrderDesc
		p.Order = &defaultOrder
	}
	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetFollowingQuery struct {
	CursorCreatedAt *time.Time   `query:"cursorCreatedAt"`
	Order           *model.Order `query:"order" validate:"omitempty,oneof=asc desc"`
}

func (p *GetFollowingQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	if p.Order == nil {
		defaultOrder := model.OrderDesc
		p.Order = &defaultOrder
	}
	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
