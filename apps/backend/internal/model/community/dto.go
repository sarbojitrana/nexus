package community

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
)

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type CreateCommunityPayload struct {
	// AdminID is set by the service from the authenticated caller, not the body.
	AdminID     string  `json:"adminId" validate:"omitempty"`
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
	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type UpdateCommunitySettingsPayload struct {
	// UserID is set by the service from the authenticated caller, not the body.
	UserID      string  `json:"userId" validate:"omitempty"`
	Name        *string `json:"name" validate:"omitempty,max=50"`
	Slug        *string `json:"slug" validate:"omitempty,max=50"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	AvatarKey   *string `json:"avatarKey"`
	BannerKey   *string `json:"bannerKey"`
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
	// TargetUserID is sourced from the :userId path param by the handler, not the body.
	TargetUserID string        `json:"memberUserId" validate:"omitempty"`
	NewRole      CommunityRole `json:"newRole" validate:"required,oneof=member moderator admin"`
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
	FollowerID  string    `json:"followerId" validate:"required,uuid"`
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

// -----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
type GetCommunitiesQuery struct {
	CursorSortValue *string      `query:"cursorSortValue"`
	CursorCreatedAt *time.Time   `query:"cursorCreatedAt"`
	Sort            *model.Sort  `query:"sort" validate:"omitempty,oneof=created_at members_count"`
	Order           *model.Order `query:"order" validate:"omitempty,oneof=asc desc"`
	Name            *string      `query:"name" validate:"omitempty"`
}

func (p *GetCommunitiesQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	if p.Sort == nil {
		defaultSort := model.SortByMembersCount
		p.Sort = &defaultSort
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

type GetCommunityMembersQuery struct {
	CursorSortValue *string        `query:"cursorSortValue"`
	CursorCreatedAt *time.Time     `query:"cursorCreatedAt"`
	Order           *model.Order   `query:"order" validate:"omitempty,oneof=asc desc"`
	Role            *CommunityRole `query:"role" validate:"omitempty,oneof=all moderator member"`
}

func (p *GetCommunityMembersQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}

	if p.Order == nil {
		defaultOrder := model.OrderAsc
		p.Order = &defaultOrder
	}

	if p.Role == nil {
		defaultRole := CommonRole
		p.Role = &defaultRole
	}

	return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type ReportCommunityPostPayload struct {
	// CommunityID is sourced from the :id path param by the handler, not the body.
	CommunityID uuid.UUID `json:"communityId" validate:"omitempty,uuid"`
	PostID      uuid.UUID `json:"postId" validate:"required,uuid"`
	Reason      string    `json:"reason" validate:"required,max=1000"`
}

func (p *ReportCommunityPostPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type ResolveCommunityPostReportPayload struct {
	// ReportID is sourced from the :reportId path param by the handler, not the body.
	ReportID      uuid.UUID             `json:"reportId" validate:"omitempty,uuid"`
	UpdatedStatus CommunityReportStatus `json:"updatedStatus" validate:"required,oneof=resolved dismissed"`
}

func (p *ResolveCommunityPostReportPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type DeleteCommunityPostPayload struct {
	PostID uuid.UUID `json:"postId" validate:"required,uuid"`
}

func (p *DeleteCommunityPostPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type BanCommunityMemberPayload struct {
	UserIDToBan string `json:"userIdToBan" validate:"required"`
}

func (p *BanCommunityMemberPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

type GetCommunityReportsQuery struct {
	// CommunityID is unused -- the repository takes it as a path-derived
	// positional argument instead, kept here only for backward-compat binding.
	CommunityID       uuid.UUID              `query:"communityId" validate:"omitempty,uuid"`
	Status            *CommunityReportStatus `query:"statusId" validate:"omitempty,oneof=pending dismissed resolved"`
	ReportedDateStart string                 `query:"reportedDateStart"`
	ReportedDateEnd   string                 `query:"reportedDateEnd"`
	CursorCreatedAt   *time.Time             `query:"cursorCreatedAt"`
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

type GetCommunityPostsQuery struct {
	CursorCreatedAt *time.Time `query:"cursorCreatedAt"`
}

func (p *GetCommunityPostsQuery) Validate() error {
	return validator.New().Struct(p)
}
