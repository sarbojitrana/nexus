package community

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
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
	Name             *string       `json:"name" validate:"omitempty,max=50"`
	NextCursor       *model.Cursor `json:"nextCursor"`
	SortBy           *string       `json:"sortBy" validate:"omitempty,oneof=created_at members_count posts_count"`
	Order            *string       `json:"order" validate:"omitempty,oneof=asc desc"`
	DateCreatedStart *string       `json:"dateCreatedStart"`
	DateCreatedEnd   *string       `json:"dateCreatedEnd"`
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

	if p.SortBy == nil {
		defaultSort := "members_count"
		p.SortBy = &defaultSort
	}

	if p.Order == nil {
		defaultOrder := "desc"
		p.Order = &defaultOrder
	}
	return nil
}

//-------------------------------------------------------------------------------------------

type ReportCommunityPostPayload struct {
	ReporterID  uuid.UUID `json:"reporterId" validate:"required,uuid"`
	CommunityID uuid.UUID `json:"communityId" validate:"required,uuid"`
	PostID      uuid.UUID `json:"postId" validate:"required,uuid"`
	Reason      string    `json:"reason" validate:"required,max=1000"`
}

func (p *ReportCommunityPostPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type ResolveCommunityPostReportPayload struct {
	ReportID     uuid.UUID             `json:"reportId" validate:"required,uuid"`
	ReportStatus CommunityReportStatus `json:"reportStatus" validate:"required,oneof=dismissed resolved pending"`
	CommunityID  uuid.UUID             `json:"communityId" validate:"required,uuid"`
	ModeratorID  uuid.UUID             `json:"moderatorId" validate:"required,uuid"`
}

func (p *ResolveCommunityPostReportPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type DeleteCommunityPostPayload struct {
	PostID      uuid.UUID `json:"postId" validate:"required,uuid"`
	CommunityID uuid.UUID `json:"communityId" validate:"required,uuid"`
	ModeratorID uuid.UUID `json:"moderatorId" validate:"required,uuid"`
}

func (p *DeleteCommunityPostPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type ChangeCommunityRolePayload struct {
	UserID      uuid.UUID     `json:"userId" validate:"required,uuid"`
	Role        CommunityRole `json:"role" validate:"required,oneof=admin moderator member"`
	CommunityID uuid.UUID     `json:"communityId" validate:"required,uuid"`
	ModeratorID uuid.UUID     `json:"moderatorId" validate:"required,uuid"`
}

func (p *ChangeCommunityRolePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type BanCommunityMemberPayload struct {
	UserID      uuid.UUID `json:"userId" validate:"required,uuid"`
	CommunityID uuid.UUID `json:"communityId" validate:"required,uuid"`
	ModeratorID uuid.UUID `json:"moderatorId" validate:"required,uuid"`
}

func (p *BanCommunityMemberPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type KickCommunityMemberPayload struct {
	UserID      uuid.UUID `json:"userId" validate:"required,uuid"`
	CommunityID uuid.UUID `json:"communityId" validate:"required,uuid"`
	ModeratorID uuid.UUID `json:"moderatorId" validate:"required,uuid"`
}

func (p *KickCommunityMemberPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type GetReportByIDPayload struct {
	ReportID uuid.UUID `json:"reportId" validate:"required,uuid"`
}

func (p *GetReportByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-------------------------------------------------------------------------------------------

type GetCommunityReportsPayload struct {
	CommunityID       uuid.UUID              `json:"communityId" validate:"required,uuid"`
	Status            *CommunityReportStatus `json:"statusId" validate:"omitempty,oneof=pending dismissed resolved"`
	ReportedDateStart string                 `json:"reportedDateStart"`
	ReportedDateEnd   string                 `json:"reportedDateEnd"`
}

func (p *GetCommunityReportsPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.Status == nil {
		defaultStatus := ReportPending
		p.Status = &defaultStatus
	}

	return nil
}
