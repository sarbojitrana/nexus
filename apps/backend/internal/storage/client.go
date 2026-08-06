package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/sarbojitrana/nexus/internal/config"
)

type Client struct {
	presign *s3.PresignClient
	bucket  string
	enabled bool
}

func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	if !cfg.AWS.IsConfigured() {
		return &Client{enabled: false}, nil
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.AWS.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.AWS.EndpointURL != "" {
			o.BaseEndpoint = aws.String(cfg.AWS.EndpointURL)
		}
		o.UsePathStyle = true
	})

	return &Client{
		presign: s3.NewPresignClient(s3Client),
		bucket:  cfg.AWS.UploadBucket,
		enabled: true,
	}, nil
}

func (c *Client) Enabled() bool {
	return c != nil && c.enabled
}

func (c *Client) PresignUpload(ctx context.Context, key, contentType string) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("storage is not configured")
	}

	req, err := c.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", fmt.Errorf("failed to presign upload: %w", err)
	}

	return req.URL, nil
}
