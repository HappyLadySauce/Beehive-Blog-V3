package logic

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	fileservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type baseLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func newBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) baseLogic {
	return baseLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func toServiceScope(scope pb.FileScope) string {
	switch scope {
	case pb.FileScope_FILE_SCOPE_AVATAR:
		return fileservice.ScopeAvatar
	case pb.FileScope_FILE_SCOPE_CONTENT_COVER:
		return fileservice.ScopeContentCover
	case pb.FileScope_FILE_SCOPE_CONTENT_IMAGE:
		return fileservice.ScopeContentImage
	case pb.FileScope_FILE_SCOPE_ATTACHMENT:
		return fileservice.ScopeAttachment
	default:
		return ""
	}
}

func toOptionalServiceScope(scope pb.FileScope) string {
	if scope == pb.FileScope_FILE_SCOPE_UNSPECIFIED {
		return ""
	}
	return toServiceScope(scope)
}

func toServiceVisibility(visibility pb.AssetVisibility) string {
	switch visibility {
	case pb.AssetVisibility_ASSET_VISIBILITY_PUBLIC:
		return fileservice.VisibilityPublic
	case pb.AssetVisibility_ASSET_VISIBILITY_PRIVATE:
		return fileservice.VisibilityPrivate
	default:
		return ""
	}
}

func toOptionalServiceVisibility(visibility pb.AssetVisibility) string {
	if visibility == pb.AssetVisibility_ASSET_VISIBILITY_UNSPECIFIED {
		return ""
	}
	return toServiceVisibility(visibility)
}

func toOptionalServiceStatus(status pb.AssetStatus) string {
	switch status {
	case pb.AssetStatus_ASSET_STATUS_PENDING:
		return fileservice.StatusPending
	case pb.AssetStatus_ASSET_STATUS_UPLOADED:
		return fileservice.StatusUploaded
	case pb.AssetStatus_ASSET_STATUS_DELETED:
		return fileservice.StatusDeleted
	default:
		return ""
	}
}

func toProtoScope(scope string) pb.FileScope {
	switch scope {
	case fileservice.ScopeAvatar:
		return pb.FileScope_FILE_SCOPE_AVATAR
	case fileservice.ScopeContentCover:
		return pb.FileScope_FILE_SCOPE_CONTENT_COVER
	case fileservice.ScopeContentImage:
		return pb.FileScope_FILE_SCOPE_CONTENT_IMAGE
	case fileservice.ScopeAttachment:
		return pb.FileScope_FILE_SCOPE_ATTACHMENT
	default:
		return pb.FileScope_FILE_SCOPE_UNSPECIFIED
	}
}

func toProtoVisibility(visibility string) pb.AssetVisibility {
	switch visibility {
	case fileservice.VisibilityPublic:
		return pb.AssetVisibility_ASSET_VISIBILITY_PUBLIC
	case fileservice.VisibilityPrivate:
		return pb.AssetVisibility_ASSET_VISIBILITY_PRIVATE
	default:
		return pb.AssetVisibility_ASSET_VISIBILITY_UNSPECIFIED
	}
}

func toProtoStatus(status string) pb.AssetStatus {
	switch status {
	case fileservice.StatusPending:
		return pb.AssetStatus_ASSET_STATUS_PENDING
	case fileservice.StatusUploaded:
		return pb.AssetStatus_ASSET_STATUS_UPLOADED
	case fileservice.StatusDeleted:
		return pb.AssetStatus_ASSET_STATUS_DELETED
	default:
		return pb.AssetStatus_ASSET_STATUS_UNSPECIFIED
	}
}

func toProtoAsset(asset *fileservice.AssetView) *pb.Asset {
	if asset == nil {
		return nil
	}
	return &pb.Asset{
		AssetId:     asset.AssetID,
		UploadId:    asset.UploadID,
		OwnerUserId: asset.OwnerUserID,
		Scope:       toProtoScope(asset.Scope),
		Visibility:  toProtoVisibility(asset.Visibility),
		Status:      toProtoStatus(asset.Status),
		Bucket:      asset.Bucket,
		ObjectKey:   asset.ObjectKey,
		PublicUrl:   asset.PublicURL,
		FileName:    asset.FileName,
		ContentType: asset.ContentType,
		ByteSize:    asset.ByteSize,
		CreatedAt:   unix(asset.CreatedAt),
		ExpiresAt:   unix(asset.ExpiresAt),
		UploadedAt:  unixPtr(asset.UploadedAt),
		DeletedAt:   unixPtr(asset.DeletedAt),
	}
}

func unix(value time.Time) int64 {
	if value.IsZero() {
		return 0
	}
	return value.Unix()
}

func unixPtr(value *time.Time) int64 {
	if value == nil || value.IsZero() {
		return 0
	}
	return value.Unix()
}

func toStatus(err error, fallback string) error {
	return errgrpcx.ToStatusWithFallback(err, errs.CodeFileInternal, fallback)
}
