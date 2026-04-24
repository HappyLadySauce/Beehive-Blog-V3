package service

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type PingService struct{ deps Dependencies }

func (s *PingService) Execute(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	if s.deps.CheckReadiness != nil {
		if err := s.deps.CheckReadiness(ctx); err != nil {
			return nil, internalErr(err)
		}
	}
	return &pb.PingResponse{Ok: true, Service: "content", Version: "v3"}, nil
}
