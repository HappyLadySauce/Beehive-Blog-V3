package storage

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

type LocalStorage struct {
	rootDir       string
	tempDir       string
	bucket        string
	uploadBaseURL string
	uploadSecret  []byte
}

func NewLocalStorage(conf config.LocalStorageConf) (*LocalStorage, error) {
	rootDir, err := filepath.Abs(strings.TrimSpace(conf.RootDir))
	if err != nil {
		return nil, fmt.Errorf("resolve local storage root: %w", err)
	}
	tempDir, err := filepath.Abs(strings.TrimSpace(conf.TempDir))
	if err != nil {
		return nil, fmt.Errorf("resolve local storage temp dir: %w", err)
	}
	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		return nil, fmt.Errorf("create local storage root: %w", err)
	}
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		return nil, fmt.Errorf("create local storage temp dir: %w", err)
	}
	uploadSecret := strings.TrimSpace(conf.UploadSecret)
	if uploadSecret == "" {
		return nil, fmt.Errorf("local upload secret is required")
	}
	return &LocalStorage{
		rootDir:       rootDir,
		tempDir:       tempDir,
		bucket:        strings.TrimSpace(conf.Bucket),
		uploadBaseURL: strings.TrimRight(strings.TrimSpace(conf.UploadBaseURL), "/"),
		uploadSecret:  []byte(uploadSecret),
	}, nil
}

func (s *LocalStorage) PresignPut(ctx context.Context, input PresignPutInput) (*PresignPutOutput, error) {
	if s == nil {
		return nil, ErrStorageDisabled
	}
	if strings.TrimSpace(input.UploadID) == "" {
		return nil, ErrStorageInvalidInput
	}
	if _, err := s.pendingPath(input.ObjectKey); err != nil {
		return nil, err
	}
	token := s.signUploadToken(input.UploadID, input.ObjectKey)
	return &PresignPutOutput{
		UploadURL: s.uploadBaseURL + "/" + url.PathEscape(input.UploadID),
		Headers: map[string]string{
			"Content-Type":    input.ContentType,
			UploadTokenHeader: token,
		},
	}, nil
}

func (s *LocalStorage) VerifyUploadToken(uploadID string, objectKey string, token string) bool {
	if s == nil || len(s.uploadSecret) == 0 {
		return false
	}
	expected := s.signUploadToken(uploadID, objectKey)
	return hmac.Equal([]byte(expected), []byte(strings.TrimSpace(token)))
}

func (s *LocalStorage) PutPending(ctx context.Context, objectKey string, body io.Reader, maxBytes int64) (*ObjectInfo, error) {
	if s == nil {
		return nil, ErrStorageDisabled
	}
	if maxBytes <= 0 {
		return nil, ErrStorageInvalidInput
	}
	target, err := s.pendingPath(objectKey)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return nil, err
	}
	tmp := target + ".uploading"
	file, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	written, copyErr := io.Copy(file, io.LimitReader(body, maxBytes+1))
	closeErr := file.Close()
	if copyErr != nil {
		_ = os.Remove(tmp)
		return nil, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return nil, closeErr
	}
	if written > maxBytes {
		_ = os.Remove(tmp)
		return nil, ErrStorageObjectTooLarge
	}
	if err := os.Rename(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return nil, err
	}
	return s.Head(ctx, s.bucket, objectKey)
}

func (s *LocalStorage) Head(ctx context.Context, bucket string, objectKey string) (*ObjectInfo, error) {
	if s == nil {
		return nil, ErrStorageDisabled
	}
	if strings.TrimSpace(bucket) != "" && bucket != s.bucket {
		return nil, ErrStorageInvalidInput
	}
	path, err := s.pendingPath(objectKey)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			path, err = s.uploadedPath(objectKey)
			if err != nil {
				return nil, err
			}
			info, err = os.Stat(path)
		}
		if err != nil {
			return nil, err
		}
	}
	if info.IsDir() {
		return nil, ErrStorageInvalidInput
	}
	return &ObjectInfo{ByteSize: info.Size()}, nil
}

func (s *LocalStorage) Commit(ctx context.Context, bucket string, objectKey string) error {
	if s == nil {
		return ErrStorageDisabled
	}
	if strings.TrimSpace(bucket) != "" && bucket != s.bucket {
		return ErrStorageInvalidInput
	}
	pending, err := s.pendingPath(objectKey)
	if err != nil {
		return err
	}
	uploaded, err := s.uploadedPath(objectKey)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(uploaded), 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(pending); errors.Is(err, os.ErrNotExist) {
		if _, uploadedErr := os.Stat(uploaded); uploadedErr == nil {
			return nil
		}
		return err
	}
	return os.Rename(pending, uploaded)
}

func (s *LocalStorage) OpenUploaded(ctx context.Context, objectKey string) (io.ReadCloser, *ObjectInfo, error) {
	if s == nil {
		return nil, nil, ErrStorageDisabled
	}
	path, err := s.uploadedPath(objectKey)
	if err != nil {
		return nil, nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, nil, err
	}
	if info.IsDir() {
		_ = file.Close()
		return nil, nil, ErrStorageInvalidInput
	}
	return file, &ObjectInfo{ByteSize: info.Size()}, nil
}

func (s *LocalStorage) Delete(ctx context.Context, bucket string, objectKey string) error {
	if s == nil {
		return ErrStorageDisabled
	}
	for _, resolve := range []func(string) (string, error){s.pendingPath, s.uploadedPath} {
		path, err := resolve(objectKey)
		if err != nil {
			return err
		}
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func (s *LocalStorage) Health(ctx context.Context) error {
	if s == nil {
		return ErrStorageDisabled
	}
	for _, dir := range []string{s.rootDir, s.tempDir} {
		probe, err := os.CreateTemp(dir, ".health-*")
		if err != nil {
			return err
		}
		name := probe.Name()
		if err := probe.Close(); err != nil {
			_ = os.Remove(name)
			return err
		}
		if err := os.Remove(name); err != nil {
			return err
		}
	}
	return nil
}

func (s *LocalStorage) pendingPath(objectKey string) (string, error) {
	return safeJoin(s.tempDir, objectKey)
}

func (s *LocalStorage) uploadedPath(objectKey string) (string, error) {
	return safeJoin(s.rootDir, objectKey)
}

func safeJoin(root string, objectKey string) (string, error) {
	normalized := filepath.Clean(filepath.FromSlash(strings.TrimSpace(objectKey)))
	if normalized == "." || normalized == "" || filepath.IsAbs(normalized) || strings.HasPrefix(normalized, ".."+string(os.PathSeparator)) || normalized == ".." {
		return "", ErrStorageInvalidInput
	}
	fullPath := filepath.Join(root, normalized)
	rel, err := filepath.Rel(root, fullPath)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", ErrStorageInvalidInput
	}
	return fullPath, nil
}

func (s *LocalStorage) signUploadToken(uploadID string, objectKey string) string {
	mac := hmac.New(sha256.New, s.uploadSecret)
	_, _ = mac.Write([]byte(strings.TrimSpace(uploadID)))
	_, _ = mac.Write([]byte{0})
	_, _ = mac.Write([]byte(strings.TrimSpace(objectKey)))
	return hex.EncodeToString(mac.Sum(nil))
}
