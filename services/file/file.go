package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/handler"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/server"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/file.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx, err := svc.NewServiceContext(c)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to initialize file service context: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := ctx.Close(); closeErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to close file service context: %v\n", closeErr)
		}
	}()
	localHTTPServer := startLocalHTTPServer(c, ctx)
	if localHTTPServer != nil {
		defer localHTTPServer.Stop()
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterFileServiceServer(grpcServer, server.NewFileServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(server.NewInternalAuthInterceptor(c).Unary())
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

func startLocalHTTPServer(c config.Config, ctx *svc.ServiceContext) *rest.Server {
	if ctx == nil || ctx.LocalStorage == nil {
		return nil
	}
	listenOn := strings.TrimSpace(c.Storage.Local.ListenOn)
	if listenOn == "" {
		return nil
	}
	restConf, err := localRESTConf(c)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "invalid local file api config: %v\n", err)
		os.Exit(1)
	}
	httpServer := rest.MustNewServer(restConf)
	handler.RegisterHandlers(httpServer, ctx)
	go func() {
		fmt.Printf("Starting local file api server at %s...\n", listenOn)
		httpServer.Start()
	}()
	return httpServer
}

func localRESTConf(c config.Config) (rest.RestConf, error) {
	host, portText, err := net.SplitHostPort(strings.TrimSpace(c.Storage.Local.ListenOn))
	if err != nil {
		return rest.RestConf{}, err
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		return rest.RestConf{}, err
	}
	return rest.RestConf{
		ServiceConf: c.ServiceConf,
		Host:        host,
		Port:        port,
		Timeout:     30000,
		MaxBytes:    maxLocalUploadBytes(c.Storage),
	}, nil
}

func maxLocalUploadBytes(c config.StorageConf) int64 {
	var maxBytes int64 = 1 << 20
	for _, value := range c.MaxBytesByScope {
		if value > maxBytes {
			maxBytes = value
		}
	}
	return maxBytes + 1024
}
