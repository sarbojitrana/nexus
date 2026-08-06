package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/notification"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
)

type NotificationHandler struct {
	Handler
	services *service.Services
}

func NewNotificationHandler(s *server.Server, services *service.Services) *NotificationHandler {
	return &NotificationHandler{Handler: NewHandler(s), services: services}
}

func (h *NotificationHandler) List(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *notification.GetNotificationsQuery) (*model.CursorPaginatedResponse[notification.Notification], error) {
		return h.services.Notification.List(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusOK, &notification.GetNotificationsQuery{})(c)
}

func (h *NotificationHandler) MarkRead(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		id, err := parseUUIDParam(c, "id")
		if err != nil {
			return err
		}
		return h.services.Notification.MarkRead(c.Request().Context(), id, middleware.GetUserID(c))
	}, http.StatusNoContent, &model.Empty{})(c)
}

func (h *NotificationHandler) MarkAllRead(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		return h.services.Notification.MarkAllRead(c.Request().Context(), middleware.GetUserID(c))
	}, http.StatusNoContent, &model.Empty{})(c)
}
