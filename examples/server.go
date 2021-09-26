package examples

import (
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"grpc/examples/article"
	pb "grpc/proto"
	"net"
)

func Run()  {
	server := grpc.NewServer()
	service := article.NewService()

	pb.RegisterArticleServiceServer(server, service)
	// 健康检查
	healthpb.RegisterHealthServer(server, service)
	lis, err := net.Listen("tcp", "127.0.0.1:8181")
	if err != nil {
		panic(err)
	}
	go func() {
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()
}