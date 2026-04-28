// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ops

import (
	"context"
	"sync"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

const readyzProbeTimeout = 2 * time.Second

type readyzProbe interface {
	Check(ctx context.Context) error
}

type readyzDependency struct {
	name  string
	probe readyzProbe
}

type readyzProbeResult struct {
	dependency string
	err        error
}

type ReadyzLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadyzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadyzLogic {
	return &ReadyzLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadyzLogic) Readyz() (resp *types.ReadyzResp, err error) {
	if l.svcCtx == nil || l.svcCtx.IdentityProbe == nil || l.svcCtx.ContentProbe == nil || l.svcCtx.FileProbe == nil {
		l.logger.Error(
			"readyz_check",
			errs.New(errs.CodeGatewayNotReady, "service is not ready"),
			logs.String("dependency", "gateway"),
			logs.String("reason", "probe_not_initialized"),
		)
		return &types.ReadyzResp{Status: "not_ready"}, errs.New(errs.CodeGatewayNotReady, "service is not ready")
	}

	probes := []readyzDependency{
		{name: "identity", probe: l.svcCtx.IdentityProbe},
		{name: "content", probe: l.svcCtx.ContentProbe},
		{name: "file", probe: l.svcCtx.FileProbe},
	}
	results := runReadyzProbes(l.ctx, probes)
	for _, result := range results {
		if result.err == nil {
			continue
		}
		l.logger.Error(
			"readyz_check",
			result.err,
			logs.String("dependency", result.dependency),
		)
		return &types.ReadyzResp{Status: "not_ready"}, errs.Wrap(result.err, errs.CodeGatewayNotReady, "service is not ready")
	}

	return &types.ReadyzResp{Status: "ready"}, nil
}

func runReadyzProbes(ctx context.Context, probes []readyzDependency) []readyzProbeResult {
	results := make(chan readyzProbeResult, len(probes))
	var wg sync.WaitGroup
	wg.Add(len(probes))
	for _, item := range probes {
		dependency := item
		go func() {
			defer wg.Done()
			probeCtx, cancel := context.WithTimeout(ctx, readyzProbeTimeout)
			defer cancel()
			results <- readyzProbeResult{
				dependency: dependency.name,
				err:        dependency.probe.Check(probeCtx),
			}
		}()
	}

	wg.Wait()
	close(results)

	collected := make([]readyzProbeResult, 0, len(probes))
	for result := range results {
		collected = append(collected, result)
	}
	return collected
}
