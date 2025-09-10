package server

import (
	"context"

	"github.com/Matyjash/Metrigo/internal/metrigo"
	pb "github.com/Matyjash/Metrigo/pb"
)

type Server struct {
	pb.UnimplementedMetrigoServer
	metrigo metrigo.Metrigo
}

func NewServer(metrigo metrigo.Metrigo) *Server {
	return &Server{
		metrigo: metrigo,
	}
}

func (s *Server) GetMemoryUsage(ctx context.Context, req *pb.MemoryUsageReq) (*pb.MemoryUsageRes, error) {
	memoryUsage, err := s.metrigo.GetMemoryUsage()
	if err != nil {
		return nil, err
	}
	return &pb.MemoryUsageRes{
		TotalB: memoryUsage.TotalB,
		UsedB:  memoryUsage.UsedB,
	}, nil
}
