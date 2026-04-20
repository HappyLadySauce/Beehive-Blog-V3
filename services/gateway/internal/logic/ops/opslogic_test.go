package ops

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	identityadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/identity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
)

type fakeProbe struct {
	checkFn func(ctx context.Context) error
}

func (f *fakeProbe) Check(ctx context.Context) error {
	if f == nil || f.checkFn == nil {
		return nil
	}

	return f.checkFn(ctx)
}

func TestHealthz(t *testing.T) {
	t.Parallel()

	logic := NewHealthzLogic(context.Background(), &svc.ServiceContext{})
	resp, err := logic.Healthz()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("expected ok, got %s", resp.Status)
	}
}

func TestReadyz(t *testing.T) {
	t.Parallel()

	logic := NewReadyzLogic(context.Background(), &svc.ServiceContext{})
	resp, err := logic.Readyz()
	if resp.Status != "not_ready" {
		t.Fatalf("expected not_ready, got %s", resp.Status)
	}
	if !errors.Is(err, errs.E(errs.CodeGatewayNotReady)) {
		t.Fatalf("expected gateway not ready error, got %v", err)
	}
}

func TestReadyzReadyWhenProbeSucceeds(t *testing.T) {
	t.Parallel()

	logic := NewReadyzLogic(context.Background(), &svc.ServiceContext{
		IdentityProbe: &fakeProbe{checkFn: func(_ context.Context) error { return nil }},
	})
	resp, err := logic.Readyz()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "ready" {
		t.Fatalf("expected ready, got %s", resp.Status)
	}
}

func TestReadyzNotReadyWhenProbeFails(t *testing.T) {
	t.Parallel()

	var probe identityadapter.ReadinessChecker = &fakeProbe{
		checkFn: func(_ context.Context) error {
			return errors.New("identity unavailable")
		},
	}

	logic := NewReadyzLogic(context.Background(), &svc.ServiceContext{IdentityProbe: probe})
	resp, err := logic.Readyz()
	if resp.Status != "not_ready" {
		t.Fatalf("expected not_ready, got %s", resp.Status)
	}
	if !errors.Is(err, errs.E(errs.CodeGatewayNotReady)) {
		t.Fatalf("expected gateway not ready error, got %v", err)
	}
}

func TestWsStub(t *testing.T) {
	t.Parallel()

	logic := NewWsLogic(context.Background(), &svc.ServiceContext{})
	resp, err := logic.Ws()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Ok {
		t.Fatalf("expected false stub response")
	}
}
