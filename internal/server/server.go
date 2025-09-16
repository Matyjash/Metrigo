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

func (s *Server) GetCpuInfo(ctx context.Context, req *pb.CpuInfoReq) (*pb.CpuInfoRes, error) {
	cpuInfo, err := s.metrigo.GetCpuInfo()
	if err != nil {
		return nil, err
	}

	var cpuInfosRes []*pb.CpuInfo
	for _, info := range cpuInfo {
		cpuInfoPb := &pb.CpuInfo{
			Id:           info.ID,
			UsagePercent: float32(info.UsagePercent),
			Frequency:    float32(info.FrequencyMhz),
		}
		cpuInfosRes = append(cpuInfosRes, cpuInfoPb)
	}

	return &pb.CpuInfoRes{CpuInfo: cpuInfosRes}, nil
}
