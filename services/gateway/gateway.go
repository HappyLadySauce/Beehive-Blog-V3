// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errhttpx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/httpx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/handler"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		requestID := ""
		if requestMeta, ok := middleware.RequestMetaFrom(ctx); ok {
			requestID = requestMeta.RequestID
		}
		if errs.Parse(err) == nil {
			err = errs.New(errs.CodeGatewayBadRequest, "bad request")
		}
		return errhttpx.BuildResponse(err, requestID)
	})

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx, err := svc.NewServiceContext(c)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to initialize gateway service context: %v\n", err)
		os.Exit(1)
	}
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
