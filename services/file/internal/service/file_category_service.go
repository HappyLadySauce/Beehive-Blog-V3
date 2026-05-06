package service

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
)

// ListFileCategories returns file categories visible to the caller.
// ListFileCategories 返回调用方可见的文件分类列表。
func (m *Manager) ListFileCategories(ctx context.Context, in ListFileCategoriesInput) ([]*FileCategoryView, error) {
	if m == nil || m.store == nil {
		return nil, serviceNotInitialized()
	}

	items, err := m.store.Categories.List(ctx, in.IncludeDisabled)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeFileInternal, "list file categories failed")
	}

	result := make([]*FileCategoryView, 0, len(items))
	for _, item := range items {
		result = append(result, toFileCategoryView(item))
	}
	return result, nil
}

// CreateFileCategory creates a file category with its allowed extensions.
// CreateFileCategory 创建文件分类及其允许后缀集合。
func (m *Manager) CreateFileCategory(ctx context.Context, in CreateFileCategoryInput) (*FileCategoryView, error) {
	if m == nil || m.store == nil {
		return nil, serviceNotInitialized()
	}

	categoryKey, err := normalizeCategoryKey(in.CategoryKey)
	if err != nil {
		return nil, err
	}
	displayName, err := normalizeDisplayName(in.DisplayName)
	if err != nil {
		return nil, err
	}
	description, err := normalizeDescription(in.Description)
	if err != nil {
		return nil, err
	}
	extensions := normalizeAllowedExtensions(in.AllowedExtensions)
	if len(extensions) == 0 {
		return nil, errs.New(errs.CodeFileInvalidExtension, "allowed_extensions is invalid")
	}

	now := time.Now().UTC()
	category := &entity.FileCategory{
		CategoryKey: categoryKey,
		DisplayName: displayName,
		Description: description,
		Enabled:     in.Enabled,
		IsDefault:   in.IsDefault,
		SortOrder:   in.SortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := m.store.Categories.Create(ctx, category, extensions); err != nil {
		if repo.IsUniqueViolation(err) {
			return nil, errs.New(errs.CodeFileInvalidArgument, "file category already exists")
		}
		return nil, errs.Wrap(err, errs.CodeFileInternal, "create file category failed")
	}
	return m.getFileCategoryView(ctx, categoryKey)
}

// UpdateFileCategory updates file category metadata without changing its key.
// UpdateFileCategory 更新文件分类元数据，不修改分类键。
func (m *Manager) UpdateFileCategory(ctx context.Context, in UpdateFileCategoryInput) (*FileCategoryView, error) {
	if m == nil || m.store == nil {
		return nil, serviceNotInitialized()
	}

	categoryKey, err := normalizeCategoryKey(in.CategoryKey)
	if err != nil {
		return nil, err
	}
	if _, err := m.store.Categories.FindByKey(ctx, categoryKey); err != nil {
		if repo.IsNotFound(err) {
			return nil, errs.New(errs.CodeFileCategoryNotFound, "file category not found")
		}
		return nil, errs.Wrap(err, errs.CodeFileInternal, "get file category failed")
	}

	displayName, err := normalizeDisplayName(in.DisplayName)
	if err != nil {
		return nil, err
	}
	description, err := normalizeDescription(in.Description)
	if err != nil {
		return nil, err
	}

	if err := m.store.Categories.Update(ctx, categoryKey, map[string]any{
		"display_name": displayName,
		"description":  description,
		"enabled":      in.Enabled,
		"sort_order":   in.SortOrder,
		"updated_at":   time.Now().UTC(),
	}); err != nil {
		return nil, errs.Wrap(err, errs.CodeFileInternal, "update file category failed")
	}
	return m.getFileCategoryView(ctx, categoryKey)
}

// UpdateFileCategoryExtensions replaces the category extension whitelist.
// UpdateFileCategoryExtensions 替换分类允许后缀白名单。
func (m *Manager) UpdateFileCategoryExtensions(ctx context.Context, in UpdateFileCategoryExtensionsInput) (*FileCategoryView, error) {
	if m == nil || m.store == nil {
		return nil, serviceNotInitialized()
	}

	categoryKey, err := normalizeCategoryKey(in.CategoryKey)
	if err != nil {
		return nil, err
	}
	if _, err := m.store.Categories.FindByKey(ctx, categoryKey); err != nil {
		if repo.IsNotFound(err) {
			return nil, errs.New(errs.CodeFileCategoryNotFound, "file category not found")
		}
		return nil, errs.Wrap(err, errs.CodeFileInternal, "get file category failed")
	}

	extensions := normalizeAllowedExtensions(in.AllowedExtensions)
	if len(extensions) == 0 {
		return nil, errs.New(errs.CodeFileInvalidExtension, "allowed_extensions is invalid")
	}
	if err := m.store.Categories.ReplaceExtensions(ctx, categoryKey, extensions); err != nil {
		return nil, errs.Wrap(err, errs.CodeFileInternal, "update file category extensions failed")
	}
	if err := m.store.Categories.Update(ctx, categoryKey, map[string]any{
		"updated_at": time.Now().UTC(),
	}); err != nil {
		return nil, errs.Wrap(err, errs.CodeFileInternal, "touch file category failed")
	}
	return m.getFileCategoryView(ctx, categoryKey)
}

// SetDefaultFileCategory marks the given category as the only default category.
// SetDefaultFileCategory 将指定分类设置为唯一默认分类。
func (m *Manager) SetDefaultFileCategory(ctx context.Context, categoryKey string) (*FileCategoryView, error) {
	if m == nil || m.store == nil {
		return nil, serviceNotInitialized()
	}

	normalizedKey, err := normalizeCategoryKey(categoryKey)
	if err != nil {
		return nil, err
	}
	category, err := m.store.Categories.FindByKey(ctx, normalizedKey)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, errs.New(errs.CodeFileCategoryNotFound, "file category not found")
		}
		return nil, errs.Wrap(err, errs.CodeFileInternal, "get file category failed")
	}
	if !category.Enabled {
		return nil, errs.New(errs.CodeFileInvalidScope, "file category is disabled")
	}
	if err := m.store.Categories.SetDefault(ctx, normalizedKey); err != nil {
		return nil, errs.Wrap(err, errs.CodeFileInternal, "set default file category failed")
	}
	return m.getFileCategoryView(ctx, normalizedKey)
}

func (m *Manager) getFileCategoryView(ctx context.Context, categoryKey string) (*FileCategoryView, error) {
	category, err := m.store.Categories.FindByKey(ctx, categoryKey)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeFileInternal, "get file category failed")
	}
	return toFileCategoryView(*category), nil
}
