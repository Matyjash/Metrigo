package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/Matyjash/Metrigo/internal/metrigo"
	"github.com/Matyjash/Metrigo/internal/server"
	"github.com/Matyjash/Metrigo/pb"
	"google.golang.org/grpc"
)

const listenPort = ":50051"

var version = "dev"

func main() {
	fmt.Printf("Metrigo version: %s\n", version)

	serverMode := flag.Bool("server", false, "Run in server mode")
	flag.Parse()

	metrigo := metrigo.NewMetrigo()
	if *serverMode {
		if err := runGrpcServer(metrigo); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Println("Running in CLI mode")

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("No command provided. Available commands: cpu, temp")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("The count of provided arguments is more than one. Trying to proceed with the first one.")
	}

	returnMessage, err := handleCommand(metrigo, args[0])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(returnMessage)
}

func runGrpcServer(metrigo metrigo.Metrigo) error {
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMetrigoServer(s, server.NewServer(metrigo))
	fmt.Printf("Server is running on port %s\n", listenPort)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

func handleCommand(metrigoMetrics metrigo.Metrigo, command string) (string, error) {
	switch command {
	case "cpu":
		cpuInfo, err := metrigoMetrics.GetCpuInfo()
		if err != nil {
			return "", err
		}
		return metrigo.CpuMessage(cpuInfo), nil
	case "temp":
		temps, err := metrigoMetrics.GetTemperatures()
		if err != nil {
			return "", err
		}
		return metrigo.TempMessage(temps), nil
	case "mem":
		memoryUsage, err := metrigoMetrics.GetMemoryUsage()
		if err != nil {
			return "", err
		}
		return metrigo.MemoryUsageMessage(memoryUsage), nil
	case "host":
		hostInfo, err := metrigoMetrics.GetHostInfo()
		if err != nil {
			return "", err
		}
		return metrigo.HostInfoMessage(hostInfo), nil
	case "net":
		netInterfaces, err := metrigoMetrics.GetNetInterfaces()
		if err != nil {
			return "", err
		}
		return metrigo.NetInterfacesMessage(netInterfaces), nil
	default:
		return "", fmt.Errorf("unknown command: %s. Available commands: cpu, temp, mem", command)
	}
}
