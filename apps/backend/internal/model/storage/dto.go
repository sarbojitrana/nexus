package storage

import "github.com/go-playground/validator/v10"

type PresignUploadPayload struct {
	MimeType string `json:"mimeType" validate:"required"`
}

func (p *PresignUploadPayload) Validate() error {
	return validator.New().Struct(p)
}

type PresignUploadResponse struct {
	UploadURL string `json:"uploadUrl"`
	Key       string `json:"key"`
}
