package follow

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type FollowCommunityPayload struct {
	CommunityID  uuid.UUID `json:"communityId" validate:"required,uuid"`
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
	FollowingID  uuid.UUID `json:"followingId" validate:"required,uuid"`
}


func (p *FollowUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type UnFollowUserPayload struct {
	FollowingID  uuid.UUID `json:"followingId" validate:"required,uuid"`
}


func (p *UnFollowUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
