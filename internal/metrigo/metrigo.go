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

	usage, err := m.metricsPuller.GetCpuUsage(true, defaultMeasureInterval)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	}

	cpuSpec, err := m.metricsPuller.GetCpusSpec()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPUs frequencies: %v", err)
	}

	if len(usage) != logicalCpuCount {
		return nil, fmt.Errorf("mismatched CPU count and usage length: %d, %d", logicalCpuCount, len(usage))
	}

	cpuInfo, err := m.buildCpuInfo(logicalCpuCount, cpuSpec, usage)
	if err != nil {
		return nil, fmt.Errorf("failed building cpu info object: %v", err)
	}

	return cpuInfo, nil
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

func (m *Metrigo) buildCpuInfo(logicalCpuCount int, cpuSpec []models.CpuSpec, usage []float64) ([]models.CpuInfo, error) {
	if len(cpuSpec) != 1 && len(cpuSpec) != logicalCpuCount {
		return nil, fmt.Errorf("not implemented yet, CPU info length (%d) and logicalCpuCount (%d) missmatch", len(cpuSpec), logicalCpuCount)
	}
	cpus := make([]models.CpuInfo, logicalCpuCount)

	for i := range logicalCpuCount {
		freq := cpuSpec[0].FrequencyMhz
		if len(cpuSpec) == logicalCpuCount {
			freq = cpuSpec[i].FrequencyMhz
		}
		cpus[i] = models.CpuInfo{
			ID:           fmt.Sprintf("cpu%d", i),
			UsagePercent: usage[i],
			CpuSpec:      models.CpuSpec{FrequencyMhz: freq},
		}
	}

	return cpus, nil
}
