package examples

import (
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"grpc/examples/article"
	pb "grpc/proto"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	Address = "127.0.0.1:8181"
)

func Run()  {
	server := grpc.NewServer()
	service := article.NewService()

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