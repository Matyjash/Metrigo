package metrigo

import (
	"fmt"
	"time"

	"github.com/Matyjash/Metrigo/internal/metrics"
	"github.com/Matyjash/Metrigo/internal/models"
)

const defaultMeasureInterval = 200 * time.Millisecond

type Metrigo struct {
	metricsPuller metrics.MetricsPuller
}

func NewMetrigo() Metrigo {
	return Metrigo{
		metricsPuller: metrics.NewGopsutilPuller(),
	}
}

func (m *Metrigo) GetCpuInfo() ([]models.CpuInfo, error) {
	logicalCpuCount, err := m.metricsPuller.GetLogicalCpuCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU count info: %v", err)
	}

	cpus := make([]models.CpuInfo, logicalCpuCount)

	usage, err := m.metricsPuller.GetCpuUsage(true, defaultMeasureInterval)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	}

	info, err := m.metricsPuller.GetCpusSpec()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPUs frequencies: %v", err)
	}

	if len(usage) != logicalCpuCount {
		return nil, fmt.Errorf("mismatched CPU count and usage length: %d, %d", logicalCpuCount, len(usage))
	}

	if len(info) == 1 {
		for i := range logicalCpuCount {
			cpus[i] = models.CpuInfo{
				ID:           fmt.Sprintf("cpu%d", i),
				UsagePercent: usage[i],
				CpuSpec: models.CpuSpec{
					FrequencyMhz: info[0].FrequencyMhz,
				},
			}
		}
	} else if len(info) == logicalCpuCount {
		for i := range logicalCpuCount {
			cpus[i] = models.CpuInfo{
				ID:           fmt.Sprintf("cpu%d", i),
				UsagePercent: usage[i],
				CpuSpec: models.CpuSpec{
					FrequencyMhz: info[i].FrequencyMhz,
				},
			}
		}
	} else {
		return nil, fmt.Errorf("not implemented yet, CPU info length (%d) and logicalCpuCount (%d) missmatch", len(info), logicalCpuCount)
	}

	return cpus, nil
}

func (m *Metrigo) GetTemperatures() ([]models.TemperatureSensor, error) {
	temps, err := m.metricsPuller.GetTemperatures()
	if err != nil {
		return nil, fmt.Errorf("failed to get temperatures: %v", err)
	}
	return temps, nil
}

func (m *Metrigo) GetMemoryUsage() (models.MemoryUsage, error) {
	usage, err := m.metricsPuller.GetVMMemoryUsage()
	if err != nil {
		return usage, fmt.Errorf("failed to get memory usage: %v", err)
	}
	return usage, nil
}
