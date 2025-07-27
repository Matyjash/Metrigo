package metrics

import (
	"fmt"
	"time"

	cpu "github.com/shirou/gopsutil/v4/cpu"
	mem "github.com/shirou/gopsutil/v4/mem"
	sensors "github.com/shirou/gopsutil/v4/sensors"
)

type MetricsPuller interface {
	GetCpuUsage(perCpu bool, interval time.Duration) ([]float64, error)
	GetPhysicalCpuCount() (int, error)
	GetLogicalCpuCount() (int, error)
	GetCpusSpec() ([]CpuSpec, error)
	GetVMMemoryUsage() (MemoryUsage, error)
	GetTemperatures() ([]TemperatureSensor, error)
}

// TODO: move data structs out of here
type CpuSpec struct {
	FrequencyMhz float64
}

type TemperatureSensor struct {
	Key   string
	Value float64
}

type MemoryUsage struct {
	UsedB  uint64
	TotalB uint64
}

type GopsutilPuller struct {
	// TODO: add logger
}

func NewGopsutilPuller() *GopsutilPuller {
	return &GopsutilPuller{}
}

func (gp *GopsutilPuller) GetCpuUsage(perCpu bool, interval time.Duration) ([]float64, error) {
	usagePercent, err := cpu.Percent(interval, perCpu)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	}
	return usagePercent, nil
}

func (gp *GopsutilPuller) GetPhysicalCpuCount() (int, error) {
	count, err := cpu.Counts(false)
	if err != nil {
		return 0, fmt.Errorf("failed to get physical CPU count: %v", err)
	}
	if count < 1 {
		return 0, fmt.Errorf("no CPUs found")
	}
	return count, nil
}

func (gp *GopsutilPuller) GetLogicalCpuCount() (int, error) {
	count, err := cpu.Counts(true)
	if err != nil {
		return 0, fmt.Errorf("failed to get logical CPU count: %v", err)
	}
	if count < 1 {
		return 0, fmt.Errorf("no CPUs found")
	}
	return count, nil
}

func (gp *GopsutilPuller) GetCpusSpec() ([]CpuSpec, error) {
	infoStats, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU stats: %v", err)
	}

	var cpusSpec []CpuSpec

	for _, cpuInfo := range infoStats {
		cpusSpec = append(cpusSpec, CpuSpec{
			FrequencyMhz: cpuInfo.Mhz,
		})
	}

	return cpusSpec, nil
}

func (gp *GopsutilPuller) GetTemperatures() ([]TemperatureSensor, error) {
	sensors, err := sensors.SensorsTemperatures()
	if err != nil {
		return []TemperatureSensor{}, fmt.Errorf("failed to get sensors temperatures: %v", err)
	}
	if len(sensors) == 0 {
		return []TemperatureSensor{}, fmt.Errorf("no temperature sensors found")
	}

	temperatureSensors := make([]TemperatureSensor, len(sensors))
	for i, sensor := range sensors {
		temperatureSensors[i] = TemperatureSensor{
			Key:   sensor.SensorKey,
			Value: sensor.Temperature,
		}
	}

	return temperatureSensors, nil
}

func (gp *GopsutilPuller) GetVMMemoryUsage() (MemoryUsage, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return MemoryUsage{}, fmt.Errorf("failed to get virtual emory stats: %v", err)
	}
	return MemoryUsage{UsedB: vmStat.Used, TotalB: vmStat.Total}, nil
}
