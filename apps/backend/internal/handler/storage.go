package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model/storage"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
)

type StorageHandler struct {
	Handler
	services *service.Services
}

func NewStorageHandler(s *server.Server, services *service.Services) *StorageHandler {
	return &StorageHandler{Handler: NewHandler(s), services: services}
}

func (h *StorageHandler) PresignUpload(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *storage.PresignUploadPayload) (*storage.PresignUploadResponse, error) {
		url, key, err := h.services.Storage.PresignUpload(c.Request().Context(), middleware.GetUserID(c), req.MimeType)
		if err != nil {
			return nil, err
		}
		return &storage.PresignUploadResponse{UploadURL: url, Key: key}, nil
	}, http.StatusOK, &storage.PresignUploadPayload{})(c)
}
