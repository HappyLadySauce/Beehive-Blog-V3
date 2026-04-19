package logic

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type RegisterLocalUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewRegisterLocalUserLogic creates a RegisterLocalUserLogic instance.
// NewRegisterLocalUserLogic 创建 RegisterLocalUserLogic 实例。
func NewRegisterLocalUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLocalUserLogic {
	return &RegisterLocalUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RegisterLocalUser registers a local account and creates the initial session.
// RegisterLocalUser 注册本地账号并创建初始会话。
func (l *RegisterLocalUserLogic) RegisterLocalUser(in *pb.RegisterLocalUserRequest) (*pb.RegisterLocalUserResponse, error) {
	// Validate local registration input before touching any dependency.
	// 在访问任何依赖前先校验本地注册输入。
	username, err := auth.NormalizeUsername(in.GetUsername())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	email, err := auth.NormalizeEmail(in.GetEmail())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	nickname, err := auth.NormalizeNickname(in.GetNickname())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := auth.ValidatePassword(in.GetPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Reject duplicated principals before entering the transaction.
	// 进入事务前先拒绝重复主体。
	if _, err := l.svcCtx.Store.Users.GetByUsername(l.ctx, username); err == nil {
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	} else if !repo.IsNotFound(err) {
		return nil, status.Errorf(codes.Internal, "lookup username failed: %v", err)
	}
	if email != "" {
		if _, err := l.svcCtx.Store.Users.GetByEmail(l.ctx, email); err == nil {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		} else if !repo.IsNotFound(err) {
			return nil, status.Errorf(codes.Internal, "lookup email failed: %v", err)
		}
	}

	hashedPassword, err := auth.HashPassword(in.GetPassword(), l.svcCtx.Config.Security.PasswordHashCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hash password failed: %v", err)
	}

	now := time.Now().UTC()
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate refresh token failed: %v", err)
	}
	refreshTokenHash := auth.HashRefreshToken(refreshToken)
	clientIP := ctxmeta.GetClientIPFromIncomingContext(l.ctx)
	authSource := auth.AuthSourceLocal

	var createdUser *entity.User
	var createdSession *entity.UserSession

	err = l.svcCtx.Store.DB().WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		store := l.svcCtx.Store.WithTx(tx)

		createdUser = &entity.User{
			Username:    username,
			Email:       stringPtr(email),
			Nickname:    stringPtr(strings.TrimSpace(nickname)),
			Role:        auth.UserRoleMember,
			Status:      auth.UserStatusActive,
			LastLoginAt: &now,
		}
		if err := store.Users.Create(l.ctx, createdUser); err != nil {
			return err
		}

		credential := &entity.CredentialLocal{
			UserID:            createdUser.ID,
			PasswordHash:      hashedPassword,
			PasswordUpdatedAt: now,
		}
		if err := store.CredentialLocals.Create(l.ctx, credential); err != nil {
			return err
		}

		createdSession = &entity.UserSession{
			UserID:     createdUser.ID,
			AuthSource: auth.AuthSourceLocal,
			Status:     auth.SessionStatusActive,
			LastSeenAt: &now,
			ExpiresAt:  now.Add(time.Duration(l.svcCtx.Config.Security.RefreshTokenTTLSeconds) * time.Second),
			IPAddress:  stringPtr(clientIP),
		}
		if err := store.UserSessions.Create(l.ctx, createdSession); err != nil {
			return err
		}

		refreshRow := &entity.RefreshToken{
			SessionID: createdSession.ID,
			TokenHash: refreshTokenHash,
			IssuedAt:  now,
			ExpiresAt: createdSession.ExpiresAt,
		}
		if err := store.RefreshTokens.Create(l.ctx, refreshRow); err != nil {
			return err
		}

		writeAudit(l.ctx, store, auditInput{
			UserID:     &createdUser.ID,
			SessionID:  &createdSession.ID,
			AuthSource: &authSource,
			EventType:  auth.AuditEventRegisterLocal,
			Result:     auth.AuditResultSuccess,
			ClientIP:   stringPtr(clientIP),
			Detail: auth.MarshalAuditDetail(map[string]any{
				"username": createdUser.Username,
			}),
		})

		return nil
	})
	if err != nil {
		if repo.IsUniqueViolation(err) {
			return nil, status.Error(codes.AlreadyExists, "username or email already exists")
		}
		return nil, status.Errorf(codes.Internal, "register local user failed: %v", err)
	}

	accessToken, accessExpiresAt, err := issueAccessToken(
		l.svcCtx.Config.Security.AccessTokenSecret,
		l.svcCtx.Config.Security.AccessTokenTTLSeconds,
		createdUser,
		createdSession,
		now,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "issue access token failed: %v", err)
	}

	l.Infof("local registration succeeded: user_id=%d username=%s", createdUser.ID, createdUser.Username)

	return &pb.RegisterLocalUserResponse{
		CurrentUser: auth.ToCurrentUser(createdUser),
		TokenPair: auth.NewTokenPair(
			accessToken,
			refreshToken,
			accessExpiresAt.Unix()-now.Unix(),
			createdSession.ID,
		),
		SessionInfo: auth.ToSessionInfo(createdSession),
	}, nil
}
