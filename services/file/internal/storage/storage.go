package storage

import (
	"context"
	"io"
	"time"
)

type PresignPutInput struct {
	UploadID    string
	Bucket      string
	ObjectKey   string
	ContentType string
	ByteSize    int64
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
	Commit(ctx context.Context, bucket string, objectKey string) error
	Delete(ctx context.Context, bucket string, objectKey string) error
	Health(ctx context.Context) error
}

type LocalObjectStorage interface {
	ObjectStorage
	PutPending(ctx context.Context, objectKey string, body io.Reader, maxBytes int64) (*ObjectInfo, error)
	OpenUploaded(ctx context.Context, objectKey string) (io.ReadCloser, *ObjectInfo, error)
	VerifyUploadToken(uploadID string, objectKey string, token string) bool
}
