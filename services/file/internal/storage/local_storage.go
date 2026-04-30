package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

type LocalStorage struct {
	rootDir string
	tempDir string
	bucket  string
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
	return &LocalStorage{
		rootDir: rootDir,
		tempDir: tempDir,
		bucket:  strings.TrimSpace(conf.Bucket),
	}, nil
}

func (s *LocalStorage) PresignPut(ctx context.Context, input PresignPutInput) (*PresignPutOutput, error) {
	if s == nil {
		return nil, ErrStorageDisabled
	}
	return nil, ErrStorageDisabled
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
