// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
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

func NewServiceContext(c config.Config) *ServiceContext {
	identityRPC := zrpc.MustNewClient(c.IdentityRPC)
	identityClient := pb.NewIdentityClient(identityRPC.Conn())
	identityProbe := identityadapter.NewProbe(identityClient)

	authMiddleware := middleware.NewAuthMiddleware(identityClient, c.Security)
	requestMetaMiddleware := middleware.NewRequestMetaMiddleware(c.Security)

	return &ServiceContext{
		Config:                c,
		IdentityRPC:           identityRPC,
		IdentityClient:        identityClient,
		IdentityProbe:         identityProbe,
		AuthMiddleware:        authMiddleware.Handle,
		RequestMetaMiddleware: requestMetaMiddleware.Handle,
	}
}
