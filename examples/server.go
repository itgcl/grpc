package examples

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"grpc/examples/article"
	pb "grpc/examples/proto"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	Address = "127.0.0.1:8181"
	AppId = "zhangsan"
	AppKey = "123456"
)

func Run()  {
	// 注册interceptor
	opts := grpc.UnaryInterceptor(interceptor)
	// 实例化server
	server := grpc.NewServer(opts)
	// 实例化article service
	service := article.NewService()
	// 注册 article
	pb.RegisterArticleServiceServer(server, service)
	// 健康检查
	healthpb.RegisterHealthServer(server, service)
	lis, err := net.Listen("tcp", Address)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Print("closed ...")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

// interceptor 拦截器
func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}
	// 继续处理请求
	return handler(ctx, req)
}

// auth 验证Token
func auth(ctx context.Context) error {
	md, has := metadata.FromIncomingContext(ctx)
	if has == false {
		return errors.New("params error")
	}
	var (
		appid  string
		appkey string
	)
	if val, ok := md["appId"]; ok {
		appid = val[0]
	}
	if val, ok := md["appKey"]; ok {
		appkey = val[0]
	}
	if appid != appid || appkey != appkey {
		return errors.New("username or password error")
	}
	return nil
}