package metrigo

import (
	"fmt"
	"strconv"

	"github.com/Matyjash/Metrigo/internal/metrics"
)

const (
	cpuMessageHeader  = "CPU metrics:\n"
	cpuMetricsMessage = "ID: %s, Usage: %s, Frequency: %s MHz"

	tempMessageHeader  = "Temperature metrics:\n"
	tempMetricsMessage = "Sensor: %s, Temperature: %s °C"

	memMessageHeader  = "Memory metrics:\n"
	memMetricsMessage = "Usage %s%%, Used: %s B, Total: %s B"
)

func CpuMessage(cpuInfo []CpuInfo) string {
	var message string = cpuMessageHeader
	for i, cpu := range cpuInfo {
		cpuID := cpu.ID
		if cpuID == "" {
			cpuID = "NA"
		}

		usagePercent := strconv.FormatFloat(cpu.UsagePercent, 'f', 2, 64)

		var frequency string = "NA"
		if cpu.FrequencyMhz != 0 {
			frequency = strconv.FormatFloat(cpu.FrequencyMhz, 'f', -1, 64)
		}

		message += fmt.Sprintf(cpuMetricsMessage, cpuID, usagePercent, frequency)
		if i != len(cpuInfo)-1 {
			message += "\n"
		}
	}
	return message
}

func TempMessage(temps []metrics.TemperatureSensor) string {
	var message string = tempMessageHeader
	for i, temp := range temps {
		sensorKey := temp.Key
		if sensorKey == "" {
			sensorKey = "NA"
		}
		value := strconv.FormatFloat(temp.Value, 'f', -1, 64)
		message += fmt.Sprintf(tempMetricsMessage, sensorKey, value)
		if i != len(temps)-1 {
			message += "\n"
		}
	}
	return message
}

func MemoryUsageMessage(memoryUsage metrics.MemoryUsage) string {
	var message string = memMessageHeader
	used := strconv.FormatUint(memoryUsage.UsedB, 10)

	total := "NA"
	if memoryUsage.TotalB != 0 {
		total = strconv.FormatUint(memoryUsage.TotalB, 10)
	}

	var usagePercent string = "NA"
	if memoryUsage.TotalB != 0 {
		usagePercent = strconv.FormatFloat((float64(memoryUsage.UsedB)/float64(memoryUsage.TotalB))*100, 'f', 2, 64)
	}

	message += fmt.Sprintf(memMetricsMessage, usagePercent, used, total)
	return message
}
