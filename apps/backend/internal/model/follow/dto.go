package follow

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type CreateFollowPayload struct {
	FollowerID  uuid.UUID `json:"followerId" validate:"required,uuid"`
	FollowingID uuid.UUID `json:"followingId" required:"required,uuid"`
}

func (p *CreateFollowPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type DeleteFollowPayload struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (p *DeleteFollowPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
