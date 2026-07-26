package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/community"
	"github.com/sarbojitrana/nexus/internal/model/post"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
)

type CommunityHandler struct {
	Handler
	services *service.Services
}

func NewCommunityHandler(s *server.Server, services *service.Services) *CommunityHandler {
	return &CommunityHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

func (h *CommunityHandler) CreateCommunity(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *community.CreateCommunityPayload) (*community.Community, error) {
		return h.services.Community.CreateCommunity(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusCreated, &community.CreateCommunityPayload{})(c)
}

func (h *CommunityHandler) GetCommunities(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *community.GetCommunitiesQuery) (*model.CursorPaginatedResponse[community.MiniCommunity], error) {
		return h.services.Community.GetCommunities(c.Request().Context(), req)
	}, http.StatusOK, &community.GetCommunitiesQuery{})(c)
}

// GetCommunityByIDOrSlug tries a UUID lookup first, falling back to slug.
func (h *CommunityHandler) GetCommunityByIDOrSlug(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *model.Empty) (*community.CommunityResponse, error) {
		idOrSlug := c.Param("idOrSlug")
		viewerID := viewerIDFromContext(c)

		if id, err := uuid.Parse(idOrSlug); err == nil {
			return h.services.Community.GetByID(c.Request().Context(), viewerID, id)
		}
		return h.services.Community.GetBySlug(c.Request().Context(), viewerID, idOrSlug)
	}, http.StatusOK, &model.Empty{})(c)
}

func (h *CommunityHandler) UpdateSettings(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *community.UpdateCommunitySettingsPayload) (*community.Community, error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Community.UpdateSettings(c.Request().Context(), middleware.GetUserID(c), communityID, req)
	}, http.StatusOK, &community.UpdateCommunitySettingsPayload{})(c)
}

func (h *CommunityHandler) DeleteCommunity(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return err
		}
		return h.services.Community.DeleteCommunity(c.Request().Context(), middleware.GetUserID(c), communityID)
	}, http.StatusNoContent, &model.Empty{})(c)
}

func (h *CommunityHandler) JoinCommunity(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *model.Empty) (*community.CommunityMember, error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Community.JoinCommunity(c.Request().Context(), middleware.GetUserID(c), communityID)
	}, http.StatusCreated, &model.Empty{})(c)
}

func (h *CommunityHandler) LeaveCommunity(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return err
		}
		return h.services.Community.LeaveCommunity(c.Request().Context(), middleware.GetUserID(c), communityID)
	}, http.StatusNoContent, &model.Empty{})(c)
}

func (h *CommunityHandler) GetMembers(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *community.GetCommunityMembersQuery) (*model.CursorPaginatedResponse[community.MiniCommunityUser], error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Community.GetMembers(c.Request().Context(), middleware.GetUserID(c), communityID, req)
	}, http.StatusOK, &community.GetCommunityMembersQuery{})(c)
}

func (h *CommunityHandler) GetCommunityPostByID(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *model.Empty) (*post.PopulatedPost, error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		postID, err := parseUUIDParam(c, "postId")
		if err != nil {
			return nil, err
		}
		return h.services.Community.GetCommunityPostByID(c.Request().Context(), middleware.GetUserID(c), communityID, postID)
	}, http.StatusOK, &model.Empty{})(c)
}

func (h *CommunityHandler) DeleteCommunityPost(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return err
		}
		postID, err := parseUUIDParam(c, "postId")
		if err != nil {
			return err
		}
		return h.services.Community.DeleteCommunityPost(c.Request().Context(), middleware.GetUserID(c), communityID, postID)
	}, http.StatusNoContent, &model.Empty{})(c)
}

func (h *CommunityHandler) ChangeMemberRole(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *community.ChangeMemberRoleInCommunityPayload) (*community.CommunityMember, error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		req.TargetUserID = c.Param("userId")
		return h.services.Community.ChangeMemberRole(c.Request().Context(), middleware.GetUserID(c), communityID, req)
	}, http.StatusOK, &community.ChangeMemberRoleInCommunityPayload{})(c)
}

func (h *CommunityHandler) BanMember(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *community.BanCommunityMemberPayload) (*community.BannedFromCommunityUser, error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Community.BanMember(c.Request().Context(), middleware.GetUserID(c), communityID, req)
	}, http.StatusCreated, &community.BanCommunityMemberPayload{})(c)
}

func (h *CommunityHandler) ReportPost(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *community.ReportCommunityPostPayload) (*community.CommunityReport, error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		req.CommunityID = communityID
		return h.services.Community.ReportPost(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusCreated, &community.ReportCommunityPostPayload{})(c)
}

func (h *CommunityHandler) GetReports(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *community.GetCommunityReportsQuery) (*model.CursorPaginatedResponse[community.CommunityReport], error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Community.GetReports(c.Request().Context(), middleware.GetUserID(c), communityID, req)
	}, http.StatusOK, &community.GetCommunityReportsQuery{})(c)
}

func (h *CommunityHandler) GetReportByID(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *model.Empty) (*community.CommunityReport, error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		reportID, err := parseUUIDParam(c, "reportId")
		if err != nil {
			return nil, err
		}
		return h.services.Community.GetReportByID(c.Request().Context(), middleware.GetUserID(c), communityID, reportID)
	}, http.StatusOK, &model.Empty{})(c)
}

func (h *CommunityHandler) ResolveReport(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *community.ResolveCommunityPostReportPayload) (*community.CommunityReport, error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		reportID, err := parseUUIDParam(c, "reportId")
		if err != nil {
			return nil, err
		}
		req.ReportID = reportID
		return h.services.Community.ResolveReport(c.Request().Context(), middleware.GetUserID(c), communityID, req)
	}, http.StatusOK, &community.ResolveCommunityPostReportPayload{})(c)
}
