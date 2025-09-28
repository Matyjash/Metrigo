package metrigo

import (
	"fmt"
	"strconv"

	"github.com/Matyjash/Metrigo/internal/models"
)

const (
	cpuMessageHeader  = "CPU metrics:\n"
	cpuMetricsMessage = "ID: %s, Usage: %s, Frequency: %s MHz"

	tempMessageHeader  = "Temperature metrics:\n"
	tempMetricsMessage = "Sensor: %s, Temperature: %s °C"

	memMessageHeader  = "Memory metrics:\n"
	memMetricsMessage = "Usage %s%%, Used: %s B, Total: %s B"

	hostMessageHeader             = "Host info:\n"
	hostMessageHostNameRow        = "Hostname: %s"
	hostMessageOSRow              = "OS: %s"
	hostMessagePlatformRow        = "Platform: %s"
	hostMessagePlatformVersionRow = "Platform version: %s"
	hostMessageUptimeRow          = "Uptime (s): %s"

	netInterfacesMessageHeader   = "Net Interfaces:\n"
	netInterfacesNameRow         = "Name: %s"
	netInterfacesIndexRow        = "Index: %d"
	netInterfacesAdressessHeader = "Adressess: \n"
	netInterfacesAdressRow       = "IP: %s"
	netInterfacesMTURow          = "MTU: %s"
)

func CpuMessage(cpuInfo []models.CpuInfo) string {
	message := cpuMessageHeader
	for i, cpu := range cpuInfo {
		cpuID := cpu.ID
		if cpuID == "" {
			cpuID = "NA"
		}

		usagePercent := strconv.FormatFloat(cpu.UsagePercent, 'f', 2, 64)

		frequency := "NA"
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

func TempMessage(temps []models.TemperatureSensor) string {
	message := tempMessageHeader
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

func MemoryUsageMessage(memoryUsage models.MemoryUsage) string {
	message := memMessageHeader
	used := strconv.FormatUint(memoryUsage.UsedB, 10)

	total := "NA"
	if memoryUsage.TotalB != 0 {
		total = strconv.FormatUint(memoryUsage.TotalB, 10)
	}

	usagePercent := "NA"
	if memoryUsage.TotalB != 0 {
		usagePercent = strconv.FormatFloat((float64(memoryUsage.UsedB)/float64(memoryUsage.TotalB))*100, 'f', 2, 64)
	}

	message += fmt.Sprintf(memMetricsMessage, usagePercent, used, total)
	return message
}

func HostInfoMessage(hostInfo models.HostInfo) string {
	message := hostMessageHeader

	hostname := "NA"
	if hostInfo.Hostname != "" {
		hostname = hostInfo.Hostname
	}
	message += fmt.Sprintf(hostMessageHostNameRow, hostname) + "\n"

	os := "NA"
	if hostInfo.OS != "" {
		os = hostInfo.OS
	}
	message += fmt.Sprintf(hostMessageOSRow, os) + "\n"

	platform := "NA"
	if hostInfo.Platform != "" {
		platform = hostInfo.Platform
	}
	message += fmt.Sprintf(hostMessagePlatformRow, platform) + "\n"

	platformVersion := "NA"
	if hostInfo.PlatformVersion != "" {
		platformVersion = hostInfo.PlatformVersion
	}
	message += fmt.Sprintf(hostMessagePlatformVersionRow, platformVersion) + "\n"

	uptime := "NA"
	if hostInfo.Uptime != 0 {
		uptime = strconv.FormatUint(hostInfo.Uptime, 10)
	}
	message += fmt.Sprintf(hostMessageUptimeRow, uptime)

	return message
}

func NetInterfacesMessage(netInferfaces []models.NetInterface) string {
	message := netInterfacesMessageHeader

	for i, iface := range netInferfaces {
		name := "NA"
		if iface.Name != "" {
			name = iface.Name
		}
		message += fmt.Sprintf(netInterfacesNameRow, name) + "\n"

		index := iface.Index
		message += fmt.Sprintf(netInterfacesIndexRow, index) + "\n"

		message += netInterfacesAdressessHeader

		for _, address := range iface.Addressess {
			adr := "NA"
			if address != "" {
				adr = address
			}
			message += fmt.Sprintf(netInterfacesAdressRow, adr) + "\n"
		}

		mtu := "NA"
		if iface.MTU != 0 {
			mtu = strconv.Itoa(iface.MTU)
		}
		message += fmt.Sprintf(netInterfacesMTURow, mtu) + "\n"

		if i != len(netInferfaces)-1 {
			message += "\n"
		}
	}

	return message
}
