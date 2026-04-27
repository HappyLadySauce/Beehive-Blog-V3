package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
)

func latestSSOFailureAuditReason(t *testing.T, deps service.Dependencies) string {
	t.Helper()

	var audit entity.IdentityAudit
	if err := deps.Store.DB().
		WithContext(context.Background()).
		Where("event_type = ? AND result = ?", auth.AuditEventFinishSSO, auth.AuditResultFailure).
		Order("id desc").
		First(&audit).Error; err != nil {
		t.Fatalf("failed to load latest sso failure audit: %v", err)
	}

	detail := map[string]any{}
	if err := json.Unmarshal(audit.Detail, &detail); err != nil {
		t.Fatalf("failed to parse audit detail: %v", err)
	}

	reason, _ := detail["reason"].(string)
	return reason
}

func writeGitHubUserResponse(w http.ResponseWriter, id int64, login, name, email string) {
	payload := map[string]any{
		"login": login,
		"name":  name,
	}
	if id > 0 {
		payload["id"] = id
	}
	if email != "" {
		payload["email"] = email
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func writeGitHubEmailsResponse(w http.ResponseWriter, payload []map[string]any) {
	if payload == nil {
		payload = []map[string]any{}
	}

	_ = json.NewEncoder(w).Encode(payload)
}
