package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/post"
	"github.com/sarbojitrana/nexus/internal/model/user"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
)

type UserHandler struct {
	Handler
	services *service.Services
}

func NewUserHandler(s *server.Server, services *service.Services) *UserHandler {
	return &UserHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

func viewerIDFromContext(c echo.Context) *string {
	if id := middleware.GetUserID(c); id != "" {
		return &id
	}
	return nil
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *user.GetUsersQuery) (*model.CursorPaginatedResponse[user.MiniUser], error) {
		return h.services.User.List(c.Request().Context(), viewerIDFromContext(c), req)
	}, http.StatusOK, &user.GetUsersQuery{})(c)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *user.GetUserByIDQuery) (*user.User, error) {
		return h.services.User.GetByID(c.Request().Context(), viewerIDFromContext(c), req.ID)
	}, http.StatusOK, &user.GetUserByIDQuery{})(c)
}

func (h *UserHandler) GetMe(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *model.Empty) (*user.User, error) {
		return h.services.User.GetOrProvisionMe(c.Request().Context(), middleware.GetUserID(c))
	}, http.StatusOK, &model.Empty{})(c)
}

func (h *UserHandler) GetMySettings(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *model.Empty) (*user.User, error) {
		return h.services.User.GetSettings(c.Request().Context(), middleware.GetUserID(c))
	}, http.StatusOK, &model.Empty{})(c)
}

func (h *UserHandler) UpdateMySettings(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *user.UpdateUserSettingsPayload) (*user.User, error) {
		return h.services.User.UpdateSettings(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusOK, &user.UpdateUserSettingsPayload{})(c)
}

func (h *UserHandler) NotifyPasswordChanged(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		return h.services.User.NotifyPasswordChanged(c.Request().Context(), middleware.GetUserID(c))
	}, http.StatusNoContent, &model.Empty{})(c)
}

func (h *UserHandler) UpdateMe(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *user.UpdateUserPayload) (*user.User, error) {
		return h.services.User.UpdateProfile(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusOK, &user.UpdateUserPayload{})(c)
}

func (h *UserHandler) DeleteMe(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		return h.services.User.DeleteAccount(c.Request().Context(), middleware.GetUserID(c))
	}, http.StatusNoContent, &model.Empty{})(c)
}

func (h *UserHandler) GetPostsByUserID(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *user.GetPostsByUserIDPayload) (*model.CursorPaginatedResponse[post.PopulatedPost], error) {
		return h.services.User.GetPostsByUser(c.Request().Context(), viewerIDFromContext(c), c.Param("id"), req)
	}, http.StatusOK, &user.GetPostsByUserIDPayload{})(c)
}
