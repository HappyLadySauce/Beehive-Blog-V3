package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client  *s3.Client
	presign *s3.PresignClient
	enabled bool
	bucket  string
}

func NewS3Storage(ctx context.Context, conf config.ObjectStorageConf) (*S3Storage, error) {
	if !conf.Enabled {
		return &S3Storage{}, nil
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
	return &S3Storage{
		client:  client,
		presign: s3.NewPresignClient(client),
		enabled: true,
		bucket:  strings.TrimSpace(conf.Bucket),
	}, nil
}

func (s *S3Storage) PresignPut(ctx context.Context, input PresignPutInput) (*PresignPutOutput, error) {
	if s == nil || !s.enabled || s.presign == nil {
		return nil, ErrStorageDisabled
	}
	result, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(input.Bucket),
		Key:         aws.String(input.ObjectKey),
		ContentType: aws.String(input.ContentType),
	}, s3.WithPresignExpires(input.Expires))
	if err != nil {
		return nil, err
	}
	return &PresignPutOutput{
		UploadURL: result.URL,
		Headers: map[string]string{
			"Content-Type": input.ContentType,
		},
	}, nil
}

func (s *S3Storage) Head(ctx context.Context, bucket string, objectKey string) (*ObjectInfo, error) {
	if s == nil || !s.enabled || s.client == nil {
		return nil, ErrStorageDisabled
	}
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}
	return &ObjectInfo{
		ByteSize:    aws.ToInt64(result.ContentLength),
		ContentType: aws.ToString(result.ContentType),
	}, nil
}

func (s *S3Storage) Delete(ctx context.Context, bucket string, objectKey string) error {
	if s == nil || !s.enabled || s.client == nil {
		return ErrStorageDisabled
	}
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	return err
}

func (s *S3Storage) Health(ctx context.Context) error {
	if s == nil || !s.enabled || s.client == nil {
		return ErrStorageDisabled
	}
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	return err
}
