package upload

import (
	"context"
	"fmt"
	"mime"
	"path"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// AvatarPresignInput describes a requested direct avatar upload.
// AvatarPresignInput 描述一次头像直传预签名请求。
type AvatarPresignInput struct {
	UserID      string
	FileName    string
	ContentType string
	ByteSize    int64
}

// AvatarPresignResult contains the upload URL and persisted public URL.
// AvatarPresignResult 包含上传地址与可持久化访问地址。
type AvatarPresignResult struct {
	UploadURL string
	PublicURL string
	ObjectKey string
	Headers   map[string]string
	ExpiresAt time.Time
	MaxBytes  int64
}

// AvatarPresigner creates S3-compatible presigned PUT URLs.
// AvatarPresigner 创建 S3 兼容的 PUT 预签名地址。
type AvatarPresigner struct {
	conf      config.ObjectStorageConf
	presign   *s3.PresignClient
	maxBytes  int64
	ttl       time.Duration
	allowlist map[string]struct{}
}

// NewAvatarPresigner builds an avatar presigner from gateway configuration.
// NewAvatarPresigner 基于 gateway 配置构建头像预签名器。
func NewAvatarPresigner(ctx context.Context, conf config.ObjectStorageConf) (*AvatarPresigner, error) {
	if !conf.Enabled {
		return &AvatarPresigner{conf: conf}, nil
	}
	if strings.TrimSpace(conf.AvatarPrefix) == "" {
		conf.AvatarPrefix = "avatars"
	}
	if conf.PresignTTLSeconds <= 0 {
		conf.PresignTTLSeconds = 300
	}
	if conf.MaxAvatarBytes <= 0 {
		conf.MaxAvatarBytes = 2 * 1024 * 1024
	}
	allowed := conf.AllowedAvatarContentTypes
	if len(allowed) == 0 {
		allowed = []string{"image/png", "image/jpeg", "image/webp", "image/avif"}
	}
	allowlist := make(map[string]struct{}, len(allowed))
	for _, item := range allowed {
		allowlist[strings.ToLower(strings.TrimSpace(item))] = struct{}{}
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(conf.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.AccessKeyID, conf.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("load object storage config: %w", err)
	}
	client := s3.NewFromConfig(awsCfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(strings.TrimRight(conf.Endpoint, "/"))
		options.UsePathStyle = true
	})

	return &AvatarPresigner{
		conf:      conf,
		presign:   s3.NewPresignClient(client),
		maxBytes:  conf.MaxAvatarBytes,
		ttl:       time.Duration(conf.PresignTTLSeconds) * time.Second,
		allowlist: allowlist,
	}, nil
}

// Presign validates the avatar file and creates a temporary upload URL.
// Presign 校验头像文件并创建临时上传地址。
func (p *AvatarPresigner) Presign(ctx context.Context, input AvatarPresignInput) (*AvatarPresignResult, error) {
	if p == nil || !p.conf.Enabled || p.presign == nil {
		return nil, errs.New(errs.CodeGatewayInternal, "avatar storage is not configured")
	}
	contentType := strings.ToLower(strings.TrimSpace(input.ContentType))
	if _, ok := p.allowlist[contentType]; !ok {
		return nil, errs.New(errs.CodeGatewayBadRequest, "avatar content_type is not allowed")
	}
	if input.ByteSize <= 0 || input.ByteSize > p.maxBytes {
		return nil, errs.New(errs.CodeGatewayBadRequest, "avatar byte_size is invalid")
	}

	objectKey := p.objectKey(input.UserID, input.FileName, contentType)
	expiresAt := time.Now().UTC().Add(p.ttl)
	result, err := p.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.conf.Bucket),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(p.ttl))
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeGatewayInternal, "create avatar upload presign failed")
	}

	return &AvatarPresignResult{
		UploadURL: result.URL,
		PublicURL: strings.TrimRight(p.conf.PublicBaseURL, "/") + "/" + objectKey,
		ObjectKey: objectKey,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
		ExpiresAt: expiresAt,
		MaxBytes:  p.maxBytes,
	}, nil
}

func (p *AvatarPresigner) objectKey(userID string, fileName string, contentType string) string {
	ext := strings.ToLower(path.Ext(fileName))
	if ext == "" {
		extensions, _ := mime.ExtensionsByType(contentType)
		if len(extensions) > 0 {
			ext = extensions[0]
		}
	}
	if ext == ".jpeg" {
		ext = ".jpg"
	}
	return strings.Trim(strings.TrimSpace(p.conf.AvatarPrefix), "/") + "/" + strings.TrimSpace(userID) + "/" + uuid.NewString() + ext
}
