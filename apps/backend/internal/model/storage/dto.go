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

type PresignDownloadsPayload struct {
	Keys []string `json:"keys" validate:"required,max=100,dive,required"`
}

func (p *PresignDownloadsPayload) Validate() error {
	return validator.New().Struct(p)
}

type PresignDownloadsResponse struct {
	URLs map[string]string `json:"urls"`
}
