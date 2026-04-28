package storage

import (
	"context"
	"time"
)

type PresignPutInput struct {
	Bucket      string
	ObjectKey   string
	ContentType string
	Expires     time.Duration
}

type PresignPutOutput struct {
	UploadURL string
	Headers   map[string]string
}

type ObjectInfo struct {
	ByteSize    int64
	ContentType string
}

type ObjectStorage interface {
	PresignPut(ctx context.Context, input PresignPutInput) (*PresignPutOutput, error)
	Head(ctx context.Context, bucket string, objectKey string) (*ObjectInfo, error)
	Delete(ctx context.Context, bucket string, objectKey string) error
	Health(ctx context.Context) error
}
