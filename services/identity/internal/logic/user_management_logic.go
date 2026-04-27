package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

type userManagementLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

func newUserManagementLogic(ctx context.Context, svcCtx *svc.ServiceContext) userManagementLogic {
	logCtx := withLogContext(ctx)
	return userManagementLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

type ListUsersLogic struct{ userManagementLogic }
type UpdateOwnProfileLogic struct{ userManagementLogic }
type ChangeOwnPasswordLogic struct{ userManagementLogic }
type UpdateUserRoleLogic struct{ userManagementLogic }
type UpdateUserStatusLogic struct{ userManagementLogic }
type ResetUserPasswordLogic struct{ userManagementLogic }
type ListIdentityAuditsLogic struct{ userManagementLogic }

// NewListUsersLogic creates a ListUsersLogic instance.
// NewListUsersLogic 创建 ListUsersLogic 实例。
func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{userManagementLogic: newUserManagementLogic(ctx, svcCtx)}
}

// NewUpdateOwnProfileLogic creates an UpdateOwnProfileLogic instance.
// NewUpdateOwnProfileLogic 创建 UpdateOwnProfileLogic 实例。
func NewUpdateOwnProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOwnProfileLogic {
	return &UpdateOwnProfileLogic{userManagementLogic: newUserManagementLogic(ctx, svcCtx)}
}

// NewChangeOwnPasswordLogic creates a ChangeOwnPasswordLogic instance.
// NewChangeOwnPasswordLogic 创建 ChangeOwnPasswordLogic 实例。
func NewChangeOwnPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeOwnPasswordLogic {
	return &ChangeOwnPasswordLogic{userManagementLogic: newUserManagementLogic(ctx, svcCtx)}
}

// NewUpdateUserRoleLogic creates an UpdateUserRoleLogic instance.
// NewUpdateUserRoleLogic 创建 UpdateUserRoleLogic 实例。
func NewUpdateUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRoleLogic {
	return &UpdateUserRoleLogic{userManagementLogic: newUserManagementLogic(ctx, svcCtx)}
}

// NewUpdateUserStatusLogic creates an UpdateUserStatusLogic instance.
// NewUpdateUserStatusLogic 创建 UpdateUserStatusLogic 实例。
func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{userManagementLogic: newUserManagementLogic(ctx, svcCtx)}
}

// NewResetUserPasswordLogic creates a ResetUserPasswordLogic instance.
// NewResetUserPasswordLogic 创建 ResetUserPasswordLogic 实例。
func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{userManagementLogic: newUserManagementLogic(ctx, svcCtx)}
}

// NewListIdentityAuditsLogic creates a ListIdentityAuditsLogic instance.
// NewListIdentityAuditsLogic 创建 ListIdentityAuditsLogic 实例。
func NewListIdentityAuditsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListIdentityAuditsLogic {
	return &ListIdentityAuditsLogic{userManagementLogic: newUserManagementLogic(ctx, svcCtx)}
}

func (l *ListUsersLogic) ListUsers(in *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	actorUserID, err := parseID("actor_user_id", in.GetActorUserId())
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.Services.Users.ListUsers(l.ctx, identityservice.ListUsersInput{
		ActorUserID: actorUserID,
		Keyword:     in.GetKeyword(),
		Role:        roleString(in.GetRole()),
		Status:      accountStatusString(in.GetStatus()),
		Page:        int(in.GetPage()),
		PageSize:    int(in.GetPageSize()),
	})
	if err != nil {
		return nil, toStatusError(err, "list users failed")
	}

	items := make([]*pb.AdminUserView, 0, len(result.Items))
	for i := range result.Items {
		items = append(items, auth.ToAdminUserView(&result.Items[i]))
	}
	return &pb.ListUsersResponse{Items: items, Total: result.Total, Page: int32(result.Page), PageSize: int32(result.PageSize)}, nil
}

func (l *UpdateOwnProfileLogic) UpdateOwnProfile(in *pb.UpdateOwnProfileRequest) (*pb.UpdateOwnProfileResponse, error) {
	userID, err := parseID("user_id", in.GetUserId())
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.Services.Users.UpdateOwnProfile(l.ctx, identityservice.UpdateOwnProfileInput{
		UserID:    userID,
		Nickname:  in.Nickname,
		AvatarURL: in.AvatarUrl,
		ClientIP:  ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "update own profile failed")
	}
	return &pb.UpdateOwnProfileResponse{CurrentUser: auth.ToCurrentUser(result.User)}, nil
}

func (l *ChangeOwnPasswordLogic) ChangeOwnPassword(in *pb.ChangeOwnPasswordRequest) (*pb.ChangeOwnPasswordResponse, error) {
	userID, err := parseID("user_id", in.GetUserId())
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.Services.Users.ChangeOwnPassword(l.ctx, identityservice.ChangeOwnPasswordInput{
		UserID:      userID,
		OldPassword: in.GetOldPassword(),
		NewPassword: in.GetNewPassword(),
		ClientIP:    ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	}); err != nil {
		return nil, toStatusError(err, "change own password failed")
	}
	return &pb.ChangeOwnPasswordResponse{Ok: true}, nil
}

func (l *UpdateUserRoleLogic) UpdateUserRole(in *pb.UpdateUserRoleRequest) (*pb.UpdateUserRoleResponse, error) {
	actorUserID, targetUserID, err := parseActorTarget(in.GetActorUserId(), in.GetTargetUserId())
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.Services.Users.UpdateUserRole(l.ctx, identityservice.UpdateUserRoleInput{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		Role:         roleString(in.GetRole()),
		ClientIP:     ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "update user role failed")
	}
	return &pb.UpdateUserRoleResponse{User: auth.ToAdminUserView(result.User)}, nil
}

func (l *UpdateUserStatusLogic) UpdateUserStatus(in *pb.UpdateUserStatusRequest) (*pb.UpdateUserStatusResponse, error) {
	actorUserID, targetUserID, err := parseActorTarget(in.GetActorUserId(), in.GetTargetUserId())
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.Services.Users.UpdateUserStatus(l.ctx, identityservice.UpdateUserStatusInput{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		Status:       accountStatusString(in.GetStatus()),
		ClientIP:     ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "update user status failed")
	}
	return &pb.UpdateUserStatusResponse{User: auth.ToAdminUserView(result.User)}, nil
}

func (l *ResetUserPasswordLogic) ResetUserPassword(in *pb.ResetUserPasswordRequest) (*pb.ResetUserPasswordResponse, error) {
	actorUserID, targetUserID, err := parseActorTarget(in.GetActorUserId(), in.GetTargetUserId())
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.Services.Users.ResetUserPassword(l.ctx, identityservice.ResetUserPasswordInput{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		NewPassword:  in.GetNewPassword(),
		ClientIP:     ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	}); err != nil {
		return nil, toStatusError(err, "reset user password failed")
	}
	return &pb.ResetUserPasswordResponse{Ok: true}, nil
}

func (l *ListIdentityAuditsLogic) ListIdentityAudits(in *pb.ListIdentityAuditsRequest) (*pb.ListIdentityAuditsResponse, error) {
	actorUserID, err := parseID("actor_user_id", in.GetActorUserId())
	if err != nil {
		return nil, err
	}
	userID, err := parseOptionalID("user_id", in.GetUserId())
	if err != nil {
		return nil, err
	}
	startedAt, err := unixSecondsPtr("started_at", in.GetStartedAt())
	if err != nil {
		return nil, err
	}
	endedAt, err := unixSecondsPtr("ended_at", in.GetEndedAt())
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.Services.Users.ListIdentityAudits(l.ctx, identityservice.ListIdentityAuditsInput{
		ActorUserID: actorUserID,
		EventType:   in.GetEventType(),
		Result:      in.GetResult(),
		UserID:      userID,
		StartedAt:   startedAt,
		EndedAt:     endedAt,
		Page:        int(in.GetPage()),
		PageSize:    int(in.GetPageSize()),
	})
	if err != nil {
		return nil, toStatusError(err, "list identity audits failed")
	}

	items := make([]*pb.IdentityAuditView, 0, len(result.Items))
	for i := range result.Items {
		items = append(items, auth.ToIdentityAuditView(&result.Items[i]))
	}
	return &pb.ListIdentityAuditsResponse{Items: items, Total: result.Total, Page: int32(result.Page), PageSize: int32(result.PageSize)}, nil
}

func parseActorTarget(actorRaw, targetRaw string) (int64, int64, error) {
	actorUserID, err := parseID("actor_user_id", actorRaw)
	if err != nil {
		return 0, 0, err
	}
	targetUserID, err := parseID("target_user_id", targetRaw)
	if err != nil {
		return 0, 0, err
	}
	return actorUserID, targetUserID, nil
}

func roleString(value pb.Role) string {
	switch value {
	case pb.Role_ROLE_ADMIN:
		return auth.UserRoleAdmin
	case pb.Role_ROLE_MEMBER:
		return auth.UserRoleMember
	default:
		return ""
	}
}

func accountStatusString(value pb.AccountStatus) string {
	switch value {
	case pb.AccountStatus_ACCOUNT_STATUS_PENDING:
		return auth.UserStatusPending
	case pb.AccountStatus_ACCOUNT_STATUS_ACTIVE:
		return auth.UserStatusActive
	case pb.AccountStatus_ACCOUNT_STATUS_DISABLED:
		return auth.UserStatusDisabled
	case pb.AccountStatus_ACCOUNT_STATUS_LOCKED:
		return auth.UserStatusLocked
	default:
		return ""
	}
}
