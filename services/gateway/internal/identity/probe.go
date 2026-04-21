package identity

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

// ReadinessChecker checks whether identity is ready to serve traffic.
// ReadinessChecker 检查 identity 是否已准备好对外提供服务。
type ReadinessChecker interface {
	Check(ctx context.Context) error
}

// Probe performs lightweight identity readiness checks through Ping RPC.
// Probe 通过 Ping RPC 执行轻量级 identity 就绪检查。
type Probe struct {
	client       pb.IdentityClient
	internalAuth ctxmeta.InternalRPCAuth
}

// NewProbe creates a Probe instance.
// NewProbe 创建 Probe 实例。
func NewProbe(client pb.IdentityClient, rpcConf config.IdentityRPCConf) *Probe {
	return &Probe{
		client: client,
		internalAuth: ctxmeta.InternalRPCAuth{
			Token:  rpcConf.InternalAuthToken,
			Caller: rpcConf.InternalCallerName,
		},
	}
}

// Check verifies the upstream identity service is ready.
// Check 验证上游 identity 服务已就绪。
func (p *Probe) Check(ctx context.Context) error {
	if p == nil || p.client == nil {
		return errs.New(errs.CodeGatewayNotReady, "identity client is not initialized")
	}

	rpcCtx := ctxmeta.BuildIdentityOutgoingContext(ctx, ctxmeta.RequestMeta{}, p.internalAuth)
	resp, err := p.client.Ping(rpcCtx, &pb.PingRequest{})
	if err != nil {
		if parsed, ok := errgrpcx.ParseStatus(err); ok {
			return errs.Wrap(parsed, errs.CodeGatewayNotReady, "identity service is not ready")
		}
		return errs.Wrap(err, errs.CodeGatewayNotReady, "identity service is not ready")
	}
	if resp == nil || !resp.GetOk() {
		return errs.New(errs.CodeGatewayNotReady, "identity service is not ready")
	}
	if resp.GetService() == "" {
		return errs.New(errs.CodeGatewayNotReady, "identity readiness response is invalid")
	}

	return nil
}
