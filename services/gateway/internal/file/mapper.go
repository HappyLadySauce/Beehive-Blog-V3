package file

import (
	filepb "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

func BuildCreateUploadRequest(actorUserID string, req *types.FileUploadCreateReq) *filepb.CreateUploadRequest {
	return &filepb.CreateUploadRequest{
		ActorUserId: actorUserID,
		Namespace:   req.Namespace,
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
		Namespace:   req.Namespace,
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
		Namespace:   asset.GetNamespace(),
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
