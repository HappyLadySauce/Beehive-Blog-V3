package file

import (
	filepb "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

func BuildCreateUploadRequest(actorUserID string, req *types.FileUploadCreateReq) *filepb.CreateUploadRequest {
	return &filepb.CreateUploadRequest{
		ActorUserId: actorUserID,
		CategoryKey: req.CategoryKey,
		FileName:    req.FileName,
		ContentType: req.ContentType,
		ByteSize:    req.ByteSize,
		Visibility:  VisibilityToProto(req.Visibility),
	}
}

func BuildCompleteUploadRequest(actorUserID string, req *types.FileUploadCompleteReq) *filepb.CompleteUploadRequest {
	return &filepb.CompleteUploadRequest{
		ActorUserId: actorUserID,
		UploadId:    req.UploadId,
	}
}

func BuildListAssetsRequest(actorUserID string, req *types.FileAssetListReq) *filepb.ListAssetsRequest {
	return &filepb.ListAssetsRequest{
		ActorUserId: actorUserID,
		CategoryKey: req.CategoryKey,
		Status:      StatusToProto(req.Status),
		Visibility:  VisibilityToProtoOptional(req.Visibility),
		OwnerUserId: req.OwnerUserId,
		Keyword:     req.Keyword,
		Page:        int32(req.Page),
		PageSize:    int32(req.PageSize),
	}
}

func BuildGetAssetRequest(actorUserID string, req *types.FileAssetIdReq) *filepb.GetAssetRequest {
	return &filepb.GetAssetRequest{
		ActorUserId: actorUserID,
		AssetId:     req.AssetId,
	}
}

func BuildDeleteAssetRequest(actorUserID string, req *types.FileAssetIdReq) *filepb.DeleteAssetRequest {
	return &filepb.DeleteAssetRequest{
		ActorUserId: actorUserID,
		AssetId:     req.AssetId,
	}
}

func BuildUpdateFileConfigRequest(req *types.FileConfigUpdateReq) *filepb.UpdateFileConfigRequest {
	return &filepb.UpdateFileConfigRequest{
		MaxUploadBytes:    req.MaxUploadBytes,
		PresignTtlSeconds: req.PresignTtlSeconds,
	}
}

func ToFileConfigView(config *filepb.FileConfig) types.FileConfigView {
	if config == nil {
		return types.FileConfigView{}
	}
	return types.FileConfigView{
		MaxUploadBytes:    config.GetMaxUploadBytes(),
		PresignTtlSeconds: config.GetPresignTtlSeconds(),
	}
}

func ToFileConfigGetResp(resp *filepb.GetFileConfigResponse) *types.FileConfigGetResp {
	if resp == nil {
		return &types.FileConfigGetResp{}
	}
	return &types.FileConfigGetResp{Config: ToFileConfigView(resp.GetConfig())}
}

func ToFileConfigUpdateResp(resp *filepb.UpdateFileConfigResponse) *types.FileConfigUpdateResp {
	if resp == nil {
		return &types.FileConfigUpdateResp{}
	}
	return &types.FileConfigUpdateResp{Config: ToFileConfigView(resp.GetConfig())}
}

func ToCreateUploadResp(resp *filepb.CreateUploadResponse) *types.FileUploadCreateResp {
	if resp == nil {
		return &types.FileUploadCreateResp{}
	}
	return &types.FileUploadCreateResp{
		Asset:     ToAssetView(resp.GetAsset()),
		UploadUrl: resp.GetUploadUrl(),
		Headers:   resp.GetHeaders(),
		ExpiresAt: resp.GetExpiresAt(),
		MaxBytes:  resp.GetMaxBytes(),
	}
}

func ToAssetResp(resp *filepb.AssetResponse) *types.FileAssetResp {
	if resp == nil {
		return &types.FileAssetResp{}
	}
	return &types.FileAssetResp{Asset: ToAssetView(resp.GetAsset())}
}

func ToAssetListResp(resp *filepb.ListAssetsResponse) *types.FileAssetListResp {
	if resp == nil {
		return &types.FileAssetListResp{}
	}
	items := make([]types.FileAssetView, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, ToAssetView(item))
	}
	return &types.FileAssetListResp{
		Items:    items,
		Total:    resp.GetTotal(),
		Page:     int(resp.GetPage()),
		PageSize: int(resp.GetPageSize()),
	}
}

func ToAssetView(asset *filepb.Asset) types.FileAssetView {
	if asset == nil {
		return types.FileAssetView{}
	}
	return types.FileAssetView{
		AssetId:     asset.GetAssetId(),
		UploadId:    asset.GetUploadId(),
		OwnerUserId: asset.GetOwnerUserId(),
		CategoryKey: asset.GetCategoryKey(),
		Visibility:  VisibilityFromProto(asset.GetVisibility()),
		Status:      StatusFromProto(asset.GetStatus()),
		Bucket:      asset.GetBucket(),
		ObjectKey:   asset.GetObjectKey(),
		PublicUrl:   asset.GetPublicUrl(),
		FileName:    asset.GetFileName(),
		ContentType: asset.GetContentType(),
		ByteSize:    asset.GetByteSize(),
		CreatedAt:   asset.GetCreatedAt(),
		ExpiresAt:   asset.GetExpiresAt(),
		UploadedAt:  asset.GetUploadedAt(),
		DeletedAt:   asset.GetDeletedAt(),
	}
}

func BuildCreateFileCategoryRequest(req *types.FileCategoryCreateReq) *filepb.CreateFileCategoryRequest {
	return &filepb.CreateFileCategoryRequest{
		CategoryKey:       req.CategoryKey,
		DisplayName:       req.DisplayName,
		Description:       req.Description,
		Enabled:           req.Enabled,
		IsDefault:         req.IsDefault,
		SortOrder:         req.SortOrder,
		AllowedExtensions: req.AllowedExtensions,
	}
}

func BuildUpdateFileCategoryRequest(req *types.FileCategoryUpdateReq) *filepb.UpdateFileCategoryRequest {
	return &filepb.UpdateFileCategoryRequest{
		CategoryKey: req.CategoryKey,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Enabled:     req.Enabled,
		SortOrder:   req.SortOrder,
	}
}

func BuildUpdateFileCategoryExtensionsRequest(req *types.FileCategoryExtensionsUpdateReq) *filepb.UpdateFileCategoryExtensionsRequest {
	return &filepb.UpdateFileCategoryExtensionsRequest{
		CategoryKey:       req.CategoryKey,
		AllowedExtensions: req.AllowedExtensions,
	}
}

func BuildSetDefaultFileCategoryRequest(req *types.FileCategoryDefaultSetReq) *filepb.SetDefaultFileCategoryRequest {
	return &filepb.SetDefaultFileCategoryRequest{
		CategoryKey: req.CategoryKey,
	}
}

func ToFileCategoryView(category *filepb.FileCategory) types.FileCategoryView {
	if category == nil {
		return types.FileCategoryView{}
	}
	return types.FileCategoryView{
		CategoryKey:       category.GetCategoryKey(),
		DisplayName:       category.GetDisplayName(),
		Description:       category.GetDescription(),
		Enabled:           category.GetEnabled(),
		IsDefault:         category.GetIsDefault(),
		SortOrder:         category.GetSortOrder(),
		AllowedExtensions: append([]string(nil), category.GetAllowedExtensions()...),
		CreatedAt:         category.GetCreatedAt(),
		UpdatedAt:         category.GetUpdatedAt(),
	}
}

func ToFileCategoryResp(resp *filepb.FileCategoryResponse) *types.FileCategoryResp {
	if resp == nil {
		return &types.FileCategoryResp{}
	}
	return &types.FileCategoryResp{Category: ToFileCategoryView(resp.GetCategory())}
}

func ToFileCategoryListResp(resp *filepb.ListFileCategoriesResponse) *types.FileCategoryListResp {
	if resp == nil {
		return &types.FileCategoryListResp{}
	}
	items := make([]types.FileCategoryView, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, ToFileCategoryView(item))
	}
	return &types.FileCategoryListResp{Items: items}
}
