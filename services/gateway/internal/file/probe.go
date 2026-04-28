package file

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	filepb "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
)

type ReadinessChecker interface {
	Check(ctx context.Context) error
}

type Probe struct {
	client       filepb.FileServiceClient
	internalAuth ctxmeta.InternalRPCAuth
}

func NewProbe(client filepb.FileServiceClient, rpcConf config.FileRPCConf) *Probe {
	return &Probe{
		client: client,
		internalAuth: ctxmeta.InternalRPCAuth{
			Token:  rpcConf.InternalAuthToken,
			Caller: rpcConf.InternalCallerName,
		},
	}
}

func (p *Probe) Check(ctx context.Context) error {
	if p == nil || p.client == nil {
		return errs.New(errs.CodeGatewayNotReady, "file client is not initialized")
	}
	rpcCtx := ctxmeta.BuildInternalOutgoingContext(ctx, ctxmeta.RequestMeta{}, p.internalAuth, ctxmeta.AuthClaims{})
	resp, err := p.client.Ping(rpcCtx, &filepb.PingRequest{})
	if err != nil {
		if parsed, ok := errgrpcx.ParseStatus(err); ok {
			return errs.Wrap(parsed, errs.CodeGatewayNotReady, "file service is not ready")
		}
		return errs.Wrap(err, errs.CodeGatewayNotReady, "file service is not ready")
	}
	if resp == nil || !resp.GetOk() || resp.GetService() == "" {
		return errs.New(errs.CodeGatewayNotReady, "file service is not ready")
	}
	return nil
}
