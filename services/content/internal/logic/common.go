package logic

import (
	"context"
	"strconv"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	contentservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func actorFromContext(ctx context.Context) (contentservice.Actor, error) {
	claims, ok := ctxmeta.TrustedAuthClaimsFromIncomingContext(ctx)
	if !ok {
		return contentservice.Actor{}, nil
	}
	userID, err := parsePositiveInt64("user_id", claims.UserID)
	if err != nil {
		return contentservice.Actor{}, err
	}
	sessionID := int64(0)
	if claims.SessionID != "" {
		sessionID, err = parsePositiveInt64("session_id", claims.SessionID)
		if err != nil {
			return contentservice.Actor{}, err
		}
	}
	return contentservice.Actor{UserID: userID, SessionID: sessionID, Role: claims.Role}, nil
}

func parsePositiveInt64(name, raw string) (int64, error) {
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		return 0, status.Errorf(codes.InvalidArgument, "%s is invalid", name)
	}
	return value, nil
}

func toStatusError(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}
	return errgrpcx.ToStatus(err, fallbackMessage)
}
