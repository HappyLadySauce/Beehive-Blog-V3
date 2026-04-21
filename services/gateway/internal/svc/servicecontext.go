// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"fmt"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	identityadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/identity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	IdentityRPC           zrpc.Client
	IdentityClient        pb.IdentityClient
	IdentityProbe         identityadapter.ReadinessChecker
	AuthMiddleware        rest.Middleware
	RequestMetaMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	identityRPC, err := zrpc.NewClient(c.IdentityRPC.RpcClientConf)
	if err != nil {
		return nil, fmt.Errorf("create identity rpc client: %w", err)
	}
	identityClient := pb.NewIdentityClient(identityRPC.Conn())
	identityProbe := identityadapter.NewProbe(identityClient, c.IdentityRPC)

	authMiddleware := middleware.NewAuthMiddleware(identityClient, c.Security, c.IdentityRPC)
	requestMetaMiddleware, err := middleware.NewRequestMetaMiddleware(c.Security)
	if err != nil {
		return nil, fmt.Errorf("create request meta middleware: %w", err)
	}

	return &ServiceContext{
		Config:                c,
		IdentityRPC:           identityRPC,
		IdentityClient:        identityClient,
		IdentityProbe:         identityProbe,
		AuthMiddleware:        authMiddleware.Handle,
		RequestMetaMiddleware: requestMetaMiddleware.Handle,
	}, nil
}
