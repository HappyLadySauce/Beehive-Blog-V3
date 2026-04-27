package identity

import (
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

func BuildListUsersRequest(actorUserID string, req *types.AdminUserListReq) *pb.ListUsersRequest {
	return &pb.ListUsersRequest{
		ActorUserId: actorUserID,
		Keyword:     req.Keyword,
		Role:        toProtoRole(req.Role),
		Status:      toProtoAccountStatus(req.Status),
		Page:        int32(req.Page),
		PageSize:    int32(req.PageSize),
	}
}

// ValidateListUsersFilters rejects unsupported HTTP query enum values before proto mapping.
// ValidateListUsersFilters 在 proto 映射前拒绝不支持的 HTTP 查询枚举值。

func ValidateListUsersFilters(req *types.AdminUserListReq) error {
	if _, err := toOptionalListRole(req.Role); err != nil {
		return err
	}
	if _, err := toOptionalListStatus(req.Status); err != nil {
		return err
	}
	return nil
}

// BuildUpdateUserRoleRequest maps studio user role updates.
// BuildUpdateUserRoleRequest 转换 Studio 用户角色更新请求。

func BuildUpdateUserRoleRequest(actorUserID string, req *types.AdminUpdateUserRoleReq) *pb.UpdateUserRoleRequest {
	return &pb.UpdateUserRoleRequest{
		ActorUserId:  actorUserID,
		TargetUserId: req.UserId,
		Role:         toProtoRole(req.Role),
	}
}

// BuildUpdateUserStatusRequest maps studio user status updates.
// BuildUpdateUserStatusRequest 转换 Studio 用户状态更新请求。

func BuildUpdateUserStatusRequest(actorUserID string, req *types.AdminUpdateUserStatusReq) *pb.UpdateUserStatusRequest {
	return &pb.UpdateUserStatusRequest{
		ActorUserId:  actorUserID,
		TargetUserId: req.UserId,
		Status:       toProtoAccountStatus(req.Status),
	}
}

// BuildResetUserPasswordRequest maps studio password reset requests.
// BuildResetUserPasswordRequest 转换 Studio 用户密码重置请求。

func BuildResetUserPasswordRequest(actorUserID string, req *types.AdminResetUserPasswordReq) *pb.ResetUserPasswordRequest {
	return &pb.ResetUserPasswordRequest{
		ActorUserId:  actorUserID,
		TargetUserId: req.UserId,
		NewPassword:  req.NewPassword,
	}
}

// BuildListAuditsRequest maps identity audit filters.
// BuildListAuditsRequest 转换 identity 审计列表过滤条件。

func BuildListAuditsRequest(actorUserID string, req *types.IdentityAuditListReq) *pb.ListIdentityAuditsRequest {
	return &pb.ListIdentityAuditsRequest{
		ActorUserId: actorUserID,
		EventType:   req.EventType,
		Result:      req.Result,
		UserId:      req.UserId,
		StartedAt:   req.StartedAt,
		EndedAt:     req.EndedAt,
		Page:        int32(req.Page),
		PageSize:    int32(req.PageSize),
	}
}

// ToRegisterResponse maps identity register response to HTTP response.
// ToRegisterResponse 将 identity 注册响应转换为 HTTP 响应。
