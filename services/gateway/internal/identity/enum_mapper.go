package identity

import (
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

func toProtoRole(role string) pb.Role {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "admin", "role_admin":
		return pb.Role_ROLE_ADMIN
	case "member", "user", "role_member":
		return pb.Role_ROLE_MEMBER
	case "guest", "role_guest":
		return pb.Role_ROLE_GUEST
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}

func toOptionalListRole(role string) (pb.Role, error) {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "":
		return pb.Role_ROLE_UNSPECIFIED, nil
	case "admin", "role_admin":
		return pb.Role_ROLE_ADMIN, nil
	case "member", "role_member":
		return pb.Role_ROLE_MEMBER, nil
	default:
		return pb.Role_ROLE_UNSPECIFIED, errs.New(errs.CodeGatewayBadRequest, "role is invalid")
	}
}

func fromProtoRole(role pb.Role) string {
	switch role {
	case pb.Role_ROLE_ADMIN:
		return "admin"
	case pb.Role_ROLE_MEMBER:
		return "member"
	case pb.Role_ROLE_GUEST:
		return "guest"
	default:
		return ""
	}
}

func toProtoAccountStatus(status string) pb.AccountStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pending", "account_status_pending":
		return pb.AccountStatus_ACCOUNT_STATUS_PENDING
	case "active", "account_status_active":
		return pb.AccountStatus_ACCOUNT_STATUS_ACTIVE
	case "disabled", "account_status_disabled":
		return pb.AccountStatus_ACCOUNT_STATUS_DISABLED
	case "locked", "account_status_locked":
		return pb.AccountStatus_ACCOUNT_STATUS_LOCKED
	default:
		return pb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED
	}
}

func toOptionalListStatus(status string) (pb.AccountStatus, error) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "":
		return pb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED, nil
	case "pending", "account_status_pending":
		return pb.AccountStatus_ACCOUNT_STATUS_PENDING, nil
	case "active", "account_status_active":
		return pb.AccountStatus_ACCOUNT_STATUS_ACTIVE, nil
	case "disabled", "account_status_disabled":
		return pb.AccountStatus_ACCOUNT_STATUS_DISABLED, nil
	case "locked", "account_status_locked":
		return pb.AccountStatus_ACCOUNT_STATUS_LOCKED, nil
	default:
		return pb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED, errs.New(errs.CodeGatewayBadRequest, "status is invalid")
	}
}

func fromProtoAccountStatus(status pb.AccountStatus) string {
	switch status {
	case pb.AccountStatus_ACCOUNT_STATUS_PENDING:
		return "pending"
	case pb.AccountStatus_ACCOUNT_STATUS_ACTIVE:
		return "active"
	case pb.AccountStatus_ACCOUNT_STATUS_DISABLED:
		return "disabled"
	case pb.AccountStatus_ACCOUNT_STATUS_LOCKED:
		return "locked"
	default:
		return ""
	}
}

func fromProtoAuthSource(authSource pb.AuthSource) string {
	switch authSource {
	case pb.AuthSource_AUTH_SOURCE_LOCAL:
		return "local"
	case pb.AuthSource_AUTH_SOURCE_SSO:
		return "sso"
	default:
		return ""
	}
}

func fromProtoSessionStatus(status pb.SessionStatus) string {
	switch status {
	case pb.SessionStatus_SESSION_STATUS_ACTIVE:
		return "active"
	case pb.SessionStatus_SESSION_STATUS_REVOKED:
		return "revoked"
	case pb.SessionStatus_SESSION_STATUS_EXPIRED:
		return "expired"
	default:
		return ""
	}
}
