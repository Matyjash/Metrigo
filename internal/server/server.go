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

	cpuInfosRes := make([]*pb.CpuInfo, len(cpuInfo))
	for i, info := range cpuInfo {
		cpuInfosRes[i] = &pb.CpuInfo{
			Id:           info.ID,
			UsagePercent: float32(info.UsagePercent),
			Frequency:    float32(info.FrequencyMhz),
		}
	}

	return &pb.CpuInfoRes{CpuInfo: cpuInfosRes}, nil
}

func (s *Server) GetTemperatures(ctx context.Context, req *pb.TemperatureReq) (*pb.TemperatureRes, error) {
	temperatures, err := s.metrigo.GetTemperatures()
	if err != nil {
		return nil, err
	}

	temperaturesRes := make([]*pb.TemperatureSensor, len(temperatures))
	for i, temperature := range temperatures {
		temperaturesRes[i] = &pb.TemperatureSensor{
			Key:   temperature.Key,
			Value: float32(temperature.Value),
		}
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

func (s *Server) GetNetInfo(ctx context.Context, req *pb.NetInfoReq) (*pb.NetInfoRes, error) {
	netInterfaces, err := s.metrigo.GetNetInterfaces()
	if err != nil {
		return nil, err
	}

	netInterfacesPb := make([]*pb.NetInterface, len(netInterfaces))
	for i, iface := range netInterfaces {
		netInterfacesPb[i] = &pb.NetInterface{
			Name:      iface.Name,
			Index:     uint32(iface.Index),
			Addresses: iface.Addressess,
			MTU:       uint64(iface.MTU),
		}
	}

	return &pb.NetInfoRes{
		Interfaces: netInterfacesPb,
	}, nil
}
