package logic

import (
	"context"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	contentservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

func TestCreateContentLogicMapsDomainError(t *testing.T) {
	t.Parallel()

	svcCtx := &svc.ServiceContext{
		Services: contentservice.NewManager(contentservice.Dependencies{}),
	}
	_, err := NewCreateContentLogic(context.Background(), svcCtx).CreateContent(&pb.CreateContentRequest{})
	if err == nil {
		t.Fatalf("expected error")
	}
	parsed, ok := errgrpcx.ParseStatus(err)
	if !ok || parsed == nil || parsed.Code != errs.CodeContentAccessForbidden {
		t.Fatalf("expected content access forbidden grpc details, got %v", parsed)
	}
}
