package article

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "grpc/examples/proto"
)

type Service struct {

}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (s *Service) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

func (s Service) CreateArticle(ctx context.Context, req *pb.RequestCreateArticle) (*pb.ReplyCreateArticle, error) {
	// 模拟数据 没有使用数据库
	err := s.CheckType(req.Type)
	if err != nil {
		return nil, err
	}
	return &pb.ReplyCreateArticle{ArticleId: 1}, nil
}

func (s Service) UpdateArticle(ctx context.Context, req *pb.RequestUpdateArticle) (*emptypb.Empty, error) {
	if err := s.checkId(req.ArticleId); err != nil {
		return nil, err
	}
	err := s.CheckType(req.Type)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s Service) DeleteArticle(ctx context.Context, req *pb.RequestDeleteArticle) (*emptypb.Empty, error) {
	if err := s.checkId(req.ArticleId); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s Service) QueryArticle(ctx context.Context, req *pb.RequestQueryArticle) (*pb.ReplyQueryArticle, error)  {
	if err := s.checkId(req.ArticleId); err != nil {
		return nil, err
	}

	return &pb.ReplyQueryArticle{
		ArticleId: 1,
		Title:     "让数字文明造福各国人民",
		Content:   "字经济是全球未来的发展方向。习主席深刻洞察人类社会发展大势，为我们积极推动数字经济和生产生活深度融合指明了前进方向，也为国际社会共同迈向数字文明新时代贡献了中国方案，必将有力推动构建人类命运共同体。",
		Author:    "张三",
		IsShow:    true,
		Type:      pb.Type_lyrics,
	}, nil
}

func (s Service) ArticleList(ctx context.Context, req *empty.Empty) (*pb.ReplyArticleList, error) {
	reply := &pb.ReplyArticleList{}
	reply.Data = append(reply.Data, &pb.ArticleItem{
		ArticleId: 1,
		Title:     "让数字文明造福各国人民",
		Content:   "字经济是全球未来的发展方向。习主席深刻洞察人类社会发展大势，为我们积极推动数字经济和生产生活深度融合指明了前进方向，也为国际社会共同迈向数字文明新时代贡献了中国方案，必将有力推动构建人类命运共同体。",
		Author:    "张三",
		IsShow:    false,
		Type:      pb.Type_lyrics,
	})
	reply.Data = append(reply.Data, &pb.ArticleItem{
		ArticleId: 2,
		Title:     "生产旺季搞拉闸限电咋回事",
		Content:   "近期，多家上市公司却发布公告称，为配合地区“能耗双控”要求限电停产。正值生产旺季，搞拉闸限电是咋回事？",
		Author:    "李四",
		IsShow:    true,
		Type:      pb.Type_novel,
	})
	return reply, nil
}
func (s Service) checkId(articleId int64) error {
	if articleId != 1 {
		return errors.New("articleId not exists")
	}
	return nil
}
func (s Service) CheckType(articleType pb.Type) error {
	switch articleType {
	case pb.Type_prose:
	case pb.Type_lyrics:
	case pb.Type_novel:
	default:
		return errors.New("article type unknown")
	}
	return nil
}