package metrics

import (
	"fmt"
	"time"

	"github.com/Matyjash/Metrigo/internal/models"
	cpu "github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	mem "github.com/shirou/gopsutil/v4/mem"
	sensors "github.com/shirou/gopsutil/v4/sensors"
)

type MetricsPuller interface {
	GetCpuUsage(perCpu bool, interval time.Duration) ([]float64, error)
	GetPhysicalCpuCount() (int, error)
	GetLogicalCpuCount() (int, error)
	GetCpusSpec() ([]models.CpuSpec, error)
	GetVMMemoryUsage() (models.MemoryUsage, error)
	GetTemperatures() ([]models.TemperatureSensor, error)
	GetHostInfo() (models.HostInfo, error)
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
		return nil, err
	}
	return usagePercent, nil
}

func (gp *GopsutilPuller) GetPhysicalCpuCount() (int, error) {
	count, err := cpu.Counts(false)
	if err != nil {
		return 0, err
	}
	if count < 1 {
		return 0, fmt.Errorf("no CPUs found")
	}
	return count, nil
}

func (gp *GopsutilPuller) GetLogicalCpuCount() (int, error) {
	count, err := cpu.Counts(true)
	if err != nil {
		return 0, err
	}
	if count < 1 {
		return 0, fmt.Errorf("no CPUs found")
	}
	return count, nil
}

func (gp *GopsutilPuller) GetCpusSpec() ([]models.CpuSpec, error) {
	infoStats, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	var cpusSpec []models.CpuSpec

	for _, cpuInfo := range infoStats {
		cpusSpec = append(cpusSpec, models.CpuSpec{
			FrequencyMhz: cpuInfo.Mhz,
		})
	}

	return cpusSpec, nil
}

func (gp *GopsutilPuller) GetTemperatures() ([]models.TemperatureSensor, error) {
	sensors, err := sensors.SensorsTemperatures()
	if err != nil {
		return []models.TemperatureSensor{}, err
	}
	if len(sensors) == 0 {
		return []models.TemperatureSensor{}, fmt.Errorf("no temperature sensors found")
	}

	temperatureSensors := make([]models.TemperatureSensor, len(sensors))
	for i, sensor := range sensors {
		temperatureSensors[i] = models.TemperatureSensor{
			Key:   sensor.SensorKey,
			Value: sensor.Temperature,
		}
	}

	return temperatureSensors, nil
}

func (gp *GopsutilPuller) GetVMMemoryUsage() (models.MemoryUsage, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return models.MemoryUsage{}, err
	}
	return models.MemoryUsage{UsedB: vmStat.Used, TotalB: vmStat.Total}, nil
}

func (gp *GopsutilPuller) GetHostInfo() (models.HostInfo, error) {
	info, err := host.Info()
	if err != nil {
		return models.HostInfo{}, err
	}
	return models.HostInfo{
		Hostname:        info.Hostname,
		OS:              info.OS,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		Uptime:          info.Uptime,
	}, nil

}
