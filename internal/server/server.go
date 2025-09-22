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

func (s *Server) GetTemperatures(ctx context.Context, req *pb.TemperatureReq) (*pb.TemperatureRes, error) {
	temperatures, err := s.metrigo.GetTemperatures()
	if err != nil {
		return nil, err
	}

	var temperaturesRes []*pb.TemperatureSensor
	for _, temperature := range temperatures {
		temperatureSensorPb := &pb.TemperatureSensor{
			Key:   temperature.Key,
			Value: float32(temperature.Value),
		}
		temperaturesRes = append(temperaturesRes, temperatureSensorPb)
	}

	return &pb.TemperatureRes{Sensors: temperaturesRes}, nil
}

func (s *Server) GetHostInfo(ctx context.Context, req *pb.HostInfoReq) (*pb.HostInfoRes, error) {
	hostInfo, err := s.metrigo.GetHostInfo()
	if err != nil {
		return nil, err
	}

	return &pb.HostInfoRes{
		Hostname:        hostInfo.Hostname,
		Os:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		Uptime:          hostInfo.Uptime,
	}, nil
}
