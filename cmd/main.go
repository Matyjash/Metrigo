package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Matyjash/Metrigo/internal/metrigo"
)

var version = "dev"

func main() {
	fmt.Printf("Metrigo version: %s\n", version)
	fmt.Println("Running in CLI mode")

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("No command provided. Available commands: cpu, temp")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("The count of provided arguments is more than one. Trying to proceed with the first one.")
	}

	metrigo := metrigo.NewMetrigo()
	returnMessage, err := handleCommand(metrigo, args[0])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(returnMessage)
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
	default:
		return "", fmt.Errorf("unknown command: %s. Available commands: cpu, temp, mem", command)
	}
}
