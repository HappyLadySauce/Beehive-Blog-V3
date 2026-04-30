package file

import filepb "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"

func VisibilityToProto(visibility string) filepb.AssetVisibility {
	switch visibility {
	case "public", "":
		return filepb.AssetVisibility_ASSET_VISIBILITY_PUBLIC
	case "private":
		return filepb.AssetVisibility_ASSET_VISIBILITY_PRIVATE
	default:
		return filepb.AssetVisibility_ASSET_VISIBILITY_UNSPECIFIED
	}
}

func VisibilityToProtoOptional(visibility string) filepb.AssetVisibility {
	switch visibility {
	case "":
		return filepb.AssetVisibility_ASSET_VISIBILITY_UNSPECIFIED
	default:
		return VisibilityToProto(visibility)
	}
}

func StatusToProto(status string) filepb.AssetStatus {
	switch status {
	case "pending":
		return filepb.AssetStatus_ASSET_STATUS_PENDING
	case "uploaded":
		return filepb.AssetStatus_ASSET_STATUS_UPLOADED
	case "deleted":
		return filepb.AssetStatus_ASSET_STATUS_DELETED
	default:
		return filepb.AssetStatus_ASSET_STATUS_UNSPECIFIED
	}
}

func VisibilityFromProto(visibility filepb.AssetVisibility) string {
	switch visibility {
	case filepb.AssetVisibility_ASSET_VISIBILITY_PUBLIC:
		return "public"
	case filepb.AssetVisibility_ASSET_VISIBILITY_PRIVATE:
		return "private"
	default:
		return ""
	}
}

func StatusFromProto(status filepb.AssetStatus) string {
	switch status {
	case filepb.AssetStatus_ASSET_STATUS_PENDING:
		return "pending"
	case filepb.AssetStatus_ASSET_STATUS_UPLOADED:
		return "uploaded"
	case filepb.AssetStatus_ASSET_STATUS_DELETED:
		return "deleted"
	default:
		return ""
	}
}
