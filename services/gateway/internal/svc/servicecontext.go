// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"fmt"

	contentpb "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	filepb "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	identityadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/identity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	IdentityClient        pb.IdentityClient
	IdentityProbe         identityadapter.ReadinessChecker
	ContentClient         contentpb.ContentClient
	ContentProbe          contentadapter.ReadinessChecker
	FileClient            filepb.FileServiceClient
	FileProbe             fileadapter.ReadinessChecker
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

	contentRPC, err := zrpc.NewClient(c.ContentRPC.RpcClientConf)
	if err != nil {
		return nil, fmt.Errorf("create content rpc client: %w", err)
	}
	contentClient := contentpb.NewContentClient(contentRPC.Conn())
	contentProbe := contentadapter.NewProbe(contentClient, c.ContentRPC)

	fileRPC, err := zrpc.NewClient(c.FileRPC.RpcClientConf)
	if err != nil {
		return nil, fmt.Errorf("create file rpc client: %w", err)
	}
	fileClient := filepb.NewFileServiceClient(fileRPC.Conn())
	fileProbe := fileadapter.NewProbe(fileClient, c.FileRPC)

	authMiddleware := middleware.NewAuthMiddleware(identityClient, c.Security, c.IdentityRPC)
	requestMetaMiddleware, err := middleware.NewRequestMetaMiddleware(c.Security)
	if err != nil {
		return nil, fmt.Errorf("create request meta middleware: %w", err)
	}

	return &ServiceContext{
		Config:                c,
		IdentityClient:        identityClient,
		IdentityProbe:         identityProbe,
		ContentClient:         contentClient,
		ContentProbe:          contentProbe,
		FileClient:            fileClient,
		FileProbe:             fileProbe,
		AuthMiddleware:        authMiddleware.Handle,
		RequestMetaMiddleware: requestMetaMiddleware.Handle,
	}, nil
}
