package identity

import (
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

func TestToLoginResponseMapsSessionEnums(t *testing.T) {
	t.Parallel()

	resp := ToLoginResponse(&pb.LoginLocalUserResponse{
		TokenPair: &pb.TokenPair{
			AccessToken:  "access",
			RefreshToken: "refresh",
			ExpiresIn:    900,
			TokenType:    "Bearer",
		},
		CurrentUser: &pb.CurrentUser{
			UserId:   "1",
			Username: "admin",
			Email:    "admin@beehive.local",
			Role:     pb.Role_ROLE_ADMIN,
			Status:   pb.AccountStatus_ACCOUNT_STATUS_ACTIVE,
		},
		SessionInfo: &pb.SessionInfo{
			SessionId:  "s1",
			UserId:     "1",
			AuthSource: pb.AuthSource_AUTH_SOURCE_LOCAL,
			Status:     pb.SessionStatus_SESSION_STATUS_ACTIVE,
		},
	})

	if resp.Session.AuthSource != "local" {
		t.Fatalf("expected local auth source, got %q", resp.Session.AuthSource)
	}
	if resp.Session.Status != "active" {
		t.Fatalf("expected active session status, got %q", resp.Session.Status)
	}
}
