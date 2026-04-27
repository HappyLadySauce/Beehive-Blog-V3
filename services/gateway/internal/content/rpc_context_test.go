package content

import (
	"context"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"google.golang.org/grpc/metadata"
)

func TestRPCContexts(t *testing.T) {
	t.Parallel()

	conf := config.ContentRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"}
	ctx := middleware.WithRequestMeta(context.Background(), ctxmeta.RequestMeta{RequestID: "req-1"})
	ctx = middleware.WithAuthContext(ctx, middleware.AuthContext{UserID: "1", SessionID: "s1", Role: "admin"})
	rpcCtx, err := StudioRPCContext(ctx, conf)
	if err != nil {
		t.Fatalf("unexpected studio context error: %v", err)
	}
	md, _ := metadata.FromOutgoingContext(rpcCtx)
	if got := md.Get(ctxmeta.MetadataKeyUserID); len(got) != 1 || got[0] != "1" {
		t.Fatalf("expected user id metadata, got %v", got)
	}
	if got := md.Get(ctxmeta.MetadataKeyInternalCaller); len(got) != 1 || got[0] != "gateway" {
		t.Fatalf("expected caller metadata, got %v", got)
	}

	publicCtx := PublicRPCContext(ctx, conf)
	publicMD, _ := metadata.FromOutgoingContext(publicCtx)
	if got := publicMD.Get(ctxmeta.MetadataKeyUserID); len(got) != 0 {
		t.Fatalf("did not expect public user metadata, got %v", got)
	}
	if got := publicMD.Get(ctxmeta.MetadataKeyInternalAuthToken); len(got) != 1 || got[0] != "secret" {
		t.Fatalf("expected public internal auth metadata, got %v", got)
	}
}
