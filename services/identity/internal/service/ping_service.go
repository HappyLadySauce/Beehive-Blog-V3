package service

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
)

// PingService executes lightweight readiness checks for identity.
// PingService 执行 identity 的轻量级就绪检查。
type PingService struct {
	deps Dependencies
}

// NewPingService creates a PingService instance.
// NewPingService 创建 PingService 实例。
func NewPingService(deps Dependencies) *PingService {
	return &PingService{deps: deps}
}

// Execute verifies identity can safely serve traffic.
// Execute 验证 identity 是否可以安全承载流量。
func (s *PingService) Execute(ctx context.Context) (*PingResult, error) {
	if s.deps.CheckReadiness == nil {
		return nil, errs.New(errs.CodeIdentityDependencyUnavailable, "identity readiness checker is not configured")
	}

	probeCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := s.deps.CheckReadiness(probeCtx); err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityDependencyUnavailable, "identity dependencies are unavailable")
	}

	return &PingResult{
		OK:      true,
		Service: "identity",
		Version: "v3",
	}, nil
}
