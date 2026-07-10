package community

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
)

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type CreateCommunityPayload struct {
	AdminID     uuid.UUID       `json:"adminId" validate:"required,uuid"`
	Name        string          `json:"name" validate:"required,max=50"`
	Slug        string          `json:"slug" validate:"required,max=50"`
	Description *string         `json:"description" validate:"omitempty,max=1000"`
	AvatarKey   *string         `json:"avatarKey"`
	BannerKey   *string         `json:"bannerKey"`
	CanPost     *PostPermission `json:"canPost"`
}

func (p *CreateCommunityPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	if p.CanPost == nil {
		defaultCanPost := AllPost
		p.CanPost = &defaultCanPost
	}
	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type UpdateCommunitySettingsPayload struct {
	UserID      uuid.UUID `jso:"userId" validate:"required,uuid"`
	Name        *string   `json:"name" validate:"omitempty,max=50"`
	Slug        *string   `json:"slug" validate:"omitempty,max=50"`
	Description *string   `json:"description" validate:"omitempty,max=1000"`
	AvatarKey   *string   `json:"avatarKey"`
	BannerKey   *string   `json:"bannerKey"`
	CanPost     *string   `json:"canPost"`
}

func (p *UpdateCommunitySettingsPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type ChangeMemberRoleInCommunityPayload struct {
	MemberUserID uuid.UUID     `json:"memberUserId" validate:"required,uuid"`
	NewRole      CommunityRole `json:"newRole" validate:"required,max=50"`
}

func (p *ChangeMemberRoleInCommunityPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetCommunitiesQuery struct {
	Name             *string    `query:"name" validate:"omitempty,max=50"`
	CursorSortValue  *string    `query:"cursorSortValue"`
	CursorCreatedAt  *time.Time `query:"cursorCreatedAt"`
	Sort             *string    `query:"sort" validate:"omitempty,oneof=created_at members_count posts_count"`
	Order            *string    `query:"order" validate:"omitempty,oneof=asc desc"`
	DateCreatedStart *string    `query:"dateCreatedStart"`
	DateCreatedEnd   *string    `query:"dateCreatedEnd"`
}

func (p *GetCommunitiesQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.Name == nil {
		defaultName := ""
		p.Name = &defaultName
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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetCommunityMembersPayload struct {
	CursorSortValue *string      `query:"cursorSortValue"`
	CursorCreatedAt *time.Time   `query:"cursorCreatedAt"`
	Sort            *model.Sort  `query:"sort" validate:"omitempty, oneof=joined_at role"`
	Order           *model.Order `query:"order" validate:"omitempty,oneof=asc desc"`
}

func (p *GetCommunityMembersPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p) ; err != nil{
		return err
	}

	if p.Sort == nil {
		defaultSort := model.SortByRole
		p.Sort = &defaultSort
	}

	if p.Order == nil {
		defaultOrder := model.OrderAsc
		p.Order = &defaultOrder
	}

	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type DeleteCommunityPostPayload struct {
	PostID      uuid.UUID `json:"postId" validate:"required,uuid"`
	CommunityID uuid.UUID `json:"communityId" validate:"required,uuid"`
	ModeratorID uuid.UUID `json:"moderatorId" validate:"required,uuid"`
}

func (p *DeleteCommunityPostPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type BanCommunityMemberPayload struct {
	UserID      uuid.UUID `json:"userId" validate:"required,uuid"`
	CommunityID uuid.UUID `json:"communityId" validate:"required,uuid"`
	ModeratorID uuid.UUID `json:"moderatorId" validate:"required,uuid"`
}

func (p *BanCommunityMemberPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type KickCommunityMemberPayload struct {
	UserID      uuid.UUID `json:"userId" validate:"required,uuid"`
	CommunityID uuid.UUID `json:"communityId" validate:"required,uuid"`
	ModeratorID uuid.UUID `json:"moderatorId" validate:"required,uuid"`
}

func (p *KickCommunityMemberPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetReportByIDQuery struct {
	ReportID uuid.UUID `query:"reportId" validate:"required,uuid"`
}

func (p *GetReportByIDQuery) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetCommunityReportsQuery struct {
	CommunityID       uuid.UUID              `query:"communityId" validate:"required,uuid"`
	Status            *CommunityReportStatus `query:"statusId" validate:"omitempty,oneof=pending dismissed resolved"`
	ReportedDateStart string                 `query:"reportedDateStart"`
	ReportedDateEnd   string                 `query:"reportedDateEnd"`
}

func (p *GetCommunityReportsQuery) Validate() error {
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
