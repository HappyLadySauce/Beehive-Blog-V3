package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/server"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/content.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx, err := svc.NewServiceContext(c)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to initialize content service context: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := ctx.Close(); closeErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to close content service context: %v\n", closeErr)
		}
	}()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterContentServer(grpcServer, server.NewContentServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(server.NewInternalAuthInterceptor(c).Unary())
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
