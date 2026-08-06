package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/storage"
)

type StorageService struct {
	server  *server.Server
	storage *storage.Client
}

func NewStorageService(s *server.Server, storage *storage.Client) *StorageService {
	return &StorageService{server: s, storage: storage}
}

func (s *StorageService) PresignUpload(ctx context.Context, userID, mimeType string) (url string, key string, err error) {
	logger := middleware.GetLoggerFromContext(ctx)

	if !s.storage.Enabled() {
		code := "STORAGE_NOT_CONFIGURED"
		return "", "", errs.NewBadRequestError("uploads are not available right now", false, &code, nil, nil)
	}

	key = fmt.Sprintf("uploads/%s/%s", userID, uuid.NewString())

	url, err = s.storage.PresignUpload(ctx, key, mimeType)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to presign upload")
		return "", "", err
	}

	logger.Info().Str("event", "upload_presigned").Str("user_id", userID).Str("key", key).Msg("upload URL presigned")
	return url, key, nil
}
