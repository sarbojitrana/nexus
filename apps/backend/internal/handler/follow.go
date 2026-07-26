package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/follow"
	"github.com/sarbojitrana/nexus/internal/model/user"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
)

type FollowHandler struct {
	Handler
	services *service.Services
}

func NewFollowHandler(s *server.Server, services *service.Services) *FollowHandler {
	return &FollowHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

func (h *FollowHandler) FollowUser(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *follow.FollowUserPayload) (*follow.UserFollow, error) {
		req.FollowingID = c.Param("id")
		return h.services.Follow.FollowUser(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusCreated, &follow.FollowUserPayload{})(c)
}

func (h *FollowHandler) UnFollowUser(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *follow.UnFollowUserPayload) error {
		req.FollowingID = c.Param("id")
		return h.services.Follow.UnFollowUser(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusNoContent, &follow.UnFollowUserPayload{})(c)
}

// IsFollowingUser mirrors GitHub's "check if following" convention: 204 if
// following, 404 if not.
func (h *FollowHandler) IsFollowingUser(c echo.Context) error {
	following, err := h.services.Follow.IsFollowingUser(c.Request().Context(), middleware.GetUserID(c), c.Param("id"))
	if err != nil {
		return err
	}
	if !*following {
		return errs.NewNotFoundError("not following", false, nil)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *FollowHandler) GetFollowers(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *follow.GetFollowersQuery) (*model.CursorPaginatedResponse[user.MiniUser], error) {
		return h.services.Follow.GetFollowers(c.Request().Context(), viewerIDFromContext(c), c.Param("id"), req)
	}, http.StatusOK, &follow.GetFollowersQuery{})(c)
}

func (h *FollowHandler) GetFollowing(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *follow.GetFollowersQuery) (*model.CursorPaginatedResponse[user.MiniUser], error) {
		return h.services.Follow.GetFollowing(c.Request().Context(), viewerIDFromContext(c), c.Param("id"), req)
	}, http.StatusOK, &follow.GetFollowersQuery{})(c)
}

func (h *FollowHandler) FollowCommunity(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *follow.FollowCommunityPayload) (*follow.CommunityFollow, error) {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		req.CommunityID = communityID
		return h.services.Follow.FollowCommunity(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusCreated, &follow.FollowCommunityPayload{})(c)
}

func (h *FollowHandler) UnFollowCommunity(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *follow.UnFollowCommunityPayload) error {
		communityID, err := parseUUIDParam(c, "id")
		if err != nil {
			return err
		}
		req.CommunityID = communityID
		return h.services.Follow.UnFollowCommunity(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusNoContent, &follow.UnFollowCommunityPayload{})(c)
}

func (h *FollowHandler) IsFollowingCommunity(c echo.Context) error {
	communityID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	following, err := h.services.Follow.IsFollowingCommunity(c.Request().Context(), middleware.GetUserID(c), communityID)
	if err != nil {
		return err
	}
	if !*following {
		return errs.NewNotFoundError("not following", false, nil)
	}
	return c.NoContent(http.StatusNoContent)
}
