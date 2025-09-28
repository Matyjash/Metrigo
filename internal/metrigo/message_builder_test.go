package metrigo

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Matyjash/Metrigo/internal/models"
)

func Test_CpuMessage(t *testing.T) {
	tests := []struct {
		name               string
		cpuInfo            []models.CpuInfo
		wantReturnContains []string
	}{
		{
			name: "formats CPU info correctly for one CPU",
			cpuInfo: []models.CpuInfo{
				{ID: "cpu0", UsagePercent: 15.5, CpuSpec: models.CpuSpec{FrequencyMhz: 3200}},
			},
			wantReturnContains: []string{fmt.Sprintf(cpuMetricsMessage, "cpu0", "15.50", "3200")},
		},
		{
			name: "formats CPU info correctly for multiple CPUs",
			cpuInfo: []models.CpuInfo{
				{ID: "cpu0", UsagePercent: 10.0, CpuSpec: models.CpuSpec{FrequencyMhz: 3000}},
				{ID: "cpu1", UsagePercent: 20.0, CpuSpec: models.CpuSpec{FrequencyMhz: 3000}},
			},
			wantReturnContains: []string{
				fmt.Sprintf(cpuMetricsMessage, "cpu0", "10.00", "3000"),
				fmt.Sprintf(cpuMetricsMessage, "cpu1", "20.00", "3000"),
			},
		},
		{
			name: "handles missing CPU ID and frequency with NA",
			cpuInfo: []models.CpuInfo{
				{UsagePercent: 5.0},
			},
			wantReturnContains: []string{fmt.Sprintf(cpuMetricsMessage, "NA", "5.00", "NA")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CpuMessage(tt.cpuInfo)
			for _, substr := range tt.wantReturnContains {
				if !strings.Contains(got, substr) {
					t.Errorf("CpuMessage() = %v, want contains %v", got, substr)
				}
			}
		})
	}
}

func Test_TempMessage(t *testing.T) {
	tests := []struct {
		name               string
		temps              []models.TemperatureSensor
		wantReturnContains []string
	}{
		{
			name: "formats temperature info correctly for one sensor",
			temps: []models.TemperatureSensor{
				{Key: "sensor1", Value: 45.5},
			},
			wantReturnContains: []string{fmt.Sprintf(tempMetricsMessage, "sensor1", "45.5")},
		},
		{
			name: "formats temperature info correctly for multiple sensors",
			temps: []models.TemperatureSensor{
				{Key: "sensor1", Value: 40.0},
				{Key: "sensor2", Value: 50.0},
			},
			wantReturnContains: []string{
				fmt.Sprintf(tempMetricsMessage, "sensor1", "40"),
				fmt.Sprintf(tempMetricsMessage, "sensor2", "50"),
			},
		},
		{
			name: "handles missing sensor key with NA",
			temps: []models.TemperatureSensor{
				{Value: 30.0},
			},
			wantReturnContains: []string{fmt.Sprintf(tempMetricsMessage, "NA", "30")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TempMessage(tt.temps)
			for _, substr := range tt.wantReturnContains {
				if !strings.Contains(got, substr) {
					t.Errorf("TempMessage() = %v, want contains %v", got, substr)
				}
			}
		})
	}
}

func Test_MemoryUsageMessage(t *testing.T) {
	tests := []struct {
		name               string
		memoryUsage        models.MemoryUsage
		wantReturnContains string
	}{
		{
			name: "formats memory usage correctly",
			memoryUsage: models.MemoryUsage{
				TotalB: 8000,
				UsedB:  4000,
			},
			wantReturnContains: fmt.Sprintf(memMetricsMessage, "50.00", "4000", "8000"),
		},
		{
			name: "handles zero total memory with NA",
			memoryUsage: models.MemoryUsage{
				TotalB: 0,
				UsedB:  4000,
			},
			wantReturnContains: fmt.Sprintf(memMetricsMessage, "NA", "4000", "NA"),
		},
		{
			name: "handles zero used memory",
			memoryUsage: models.MemoryUsage{
				TotalB: 8000,
				UsedB:  0,
			},
			wantReturnContains: fmt.Sprintf(memMetricsMessage, "0.00", "0", "8000"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MemoryUsageMessage(tt.memoryUsage)
			if !strings.Contains(got, tt.wantReturnContains) {
				t.Errorf("MemoryUsageMessage() = %v, want contains %v", got, tt.wantReturnContains)
			}
		})
	}
}

func Test_HostInfoMessage(t *testing.T) {
	tests := []struct {
		name               string
		hostInfo           models.HostInfo
		wantReturnContains string
	}{
		{
			name: "returns proper host info message",
			hostInfo: models.HostInfo{
				Hostname:        "test",
				OS:              "linux",
				Platform:        "ubuntu",
				PlatformVersion: "Ubuntu 24.04.3 LTS",
				Uptime:          10,
			},
			wantReturnContains: fmt.Sprintf(hostMessageHostNameRow, "test") + "\n" +
				fmt.Sprintf(hostMessageOSRow, "linux") + "\n" +
				fmt.Sprintf(hostMessagePlatformRow, "ubuntu") + "\n" +
				fmt.Sprintf(hostMessagePlatformVersionRow, "Ubuntu 24.04.3 LTS") + "\n" +
				fmt.Sprintf(hostMessageUptimeRow, strconv.FormatUint(10, 10)),
		},
		{
			name: "replace empty hostname with NA",
			hostInfo: models.HostInfo{
				Hostname:        "",
				OS:              "linux",
				Platform:        "ubuntu",
				PlatformVersion: "Ubuntu 24.04.3 LTS",
				Uptime:          10,
			},
			wantReturnContains: fmt.Sprintf(hostMessageHostNameRow, "NA") + "\n" +
				fmt.Sprintf(hostMessageOSRow, "linux") + "\n" +
				fmt.Sprintf(hostMessagePlatformRow, "ubuntu") + "\n" +
				fmt.Sprintf(hostMessagePlatformVersionRow, "Ubuntu 24.04.3 LTS") + "\n" +
				fmt.Sprintf(hostMessageUptimeRow, strconv.FormatUint(10, 10)),
		},
		{
			name: "replace empty os with NA",
			hostInfo: models.HostInfo{
				Hostname:        "test",
				OS:              "",
				Platform:        "ubuntu",
				PlatformVersion: "Ubuntu 24.04.3 LTS",
				Uptime:          10,
			},
			wantReturnContains: fmt.Sprintf(hostMessageHostNameRow, "test") + "\n" +
				fmt.Sprintf(hostMessageOSRow, "NA") + "\n" +
				fmt.Sprintf(hostMessagePlatformRow, "ubuntu") + "\n" +
				fmt.Sprintf(hostMessagePlatformVersionRow, "Ubuntu 24.04.3 LTS") + "\n" +
				fmt.Sprintf(hostMessageUptimeRow, strconv.FormatUint(10, 10)),
		},
		{
			name: "replace empty platform with NA",
			hostInfo: models.HostInfo{
				Hostname:        "test",
				OS:              "linux",
				Platform:        "",
				PlatformVersion: "Ubuntu 24.04.3 LTS",
				Uptime:          10,
			},
			wantReturnContains: fmt.Sprintf(hostMessageHostNameRow, "test") + "\n" +
				fmt.Sprintf(hostMessageOSRow, "linux") + "\n" +
				fmt.Sprintf(hostMessagePlatformRow, "NA") + "\n" +
				fmt.Sprintf(hostMessagePlatformVersionRow, "Ubuntu 24.04.3 LTS") + "\n" +
				fmt.Sprintf(hostMessageUptimeRow, strconv.FormatUint(10, 10)),
		},
		{
			name: "replace empty platform version with NA",
			hostInfo: models.HostInfo{
				Hostname:        "test",
				OS:              "linux",
				Platform:        "ubuntu",
				PlatformVersion: "",
				Uptime:          10,
			},
			wantReturnContains: fmt.Sprintf(hostMessageHostNameRow, "test") + "\n" +
				fmt.Sprintf(hostMessageOSRow, "linux") + "\n" +
				fmt.Sprintf(hostMessagePlatformRow, "ubuntu") + "\n" +
				fmt.Sprintf(hostMessagePlatformVersionRow, "NA") + "\n" +
				fmt.Sprintf(hostMessageUptimeRow, strconv.FormatUint(10, 10)),
		},
		{
			name: "replace zero uptime with NA",
			hostInfo: models.HostInfo{
				Hostname:        "test",
				OS:              "linux",
				Platform:        "ubuntu",
				PlatformVersion: "Ubuntu 24.04.3 LTS",
				Uptime:          0,
			},
			wantReturnContains: fmt.Sprintf(hostMessageHostNameRow, "test") + "\n" +
				fmt.Sprintf(hostMessageOSRow, "linux") + "\n" +
				fmt.Sprintf(hostMessagePlatformRow, "ubuntu") + "\n" +
				fmt.Sprintf(hostMessagePlatformVersionRow, "Ubuntu 24.04.3 LTS") + "\n" +
				fmt.Sprintf(hostMessageUptimeRow, "NA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HostInfoMessage(tt.hostInfo)
			if !strings.Contains(got, tt.wantReturnContains) {
				t.Errorf("HostInfoMessage() = %v, want contains %v", got, tt.wantReturnContains)
			}
		})
	}
}

func Test_NetInterfacesMessage(t *testing.T) {
	tests := []struct {
		name               string
		netInferfaces      []models.NetInterface
		wantReturnContains string
	}{
		{
			name: "returns proper net interfaces message",
			netInferfaces: []models.NetInterface{
				{Name: "iface1", Index: 1, Addressess: []string{"ipv4", "ipv6"}, MTU: 128},
				{Name: "iface2", Index: 3, Addressess: []string{"ipv4_2", "ipv6_2"}, MTU: 64},
			},
			wantReturnContains: netInterfacesMessageHeader +
				fmt.Sprintf(netInterfacesNameRow, "iface1") + "\n" +
				fmt.Sprintf(netInterfacesIndexRow, 1) + "\n" +
				netInterfacesAdressessHeader +
				fmt.Sprintf(netInterfacesAdressRow, "ipv4") + "\n" +
				fmt.Sprintf(netInterfacesAdressRow, "ipv6") + "\n" +
				fmt.Sprintf(netInterfacesMTURow, strconv.Itoa(128)) + "\n" +
				"\n" +
				fmt.Sprintf(netInterfacesNameRow, "iface2") + "\n" +
				fmt.Sprintf(netInterfacesIndexRow, 3) + "\n" +
				netInterfacesAdressessHeader +
				fmt.Sprintf(netInterfacesAdressRow, "ipv4_2") + "\n" +
				fmt.Sprintf(netInterfacesAdressRow, "ipv6_2") + "\n" +
				fmt.Sprintf(netInterfacesMTURow, strconv.Itoa(64)) + "\n",
		},
		{
			name: "replaces empty name with NA",
			netInferfaces: []models.NetInterface{
				{Name: "", Index: 1, Addressess: []string{"ipv4", "ipv6"}, MTU: 128},
			},
			wantReturnContains: netInterfacesMessageHeader +
				fmt.Sprintf(netInterfacesNameRow, "NA") + "\n" +
				fmt.Sprintf(netInterfacesIndexRow, 1) + "\n" +
				netInterfacesAdressessHeader +
				fmt.Sprintf(netInterfacesAdressRow, "ipv4") + "\n" +
				fmt.Sprintf(netInterfacesAdressRow, "ipv6") + "\n" +
				fmt.Sprintf(netInterfacesMTURow, strconv.Itoa(128)) + "\n",
		},
		{
			name: "replaces empty address with NA",
			netInferfaces: []models.NetInterface{
				{Name: "iface1", Index: 1, Addressess: []string{"", "ipv6"}, MTU: 128},
			},
			wantReturnContains: netInterfacesMessageHeader +
				fmt.Sprintf(netInterfacesNameRow, "iface1") + "\n" +
				fmt.Sprintf(netInterfacesIndexRow, 1) + "\n" +
				netInterfacesAdressessHeader +
				fmt.Sprintf(netInterfacesAdressRow, "NA") + "\n" +
				fmt.Sprintf(netInterfacesAdressRow, "ipv6") + "\n" +
				fmt.Sprintf(netInterfacesMTURow, strconv.Itoa(128)) + "\n",
		},
		{
			name: "replaces zero MTU with NA",
			netInferfaces: []models.NetInterface{
				{Name: "iface1", Index: 1, Addressess: []string{"ipv4", "ipv6"}, MTU: 0},
			},
			wantReturnContains: netInterfacesMessageHeader +
				fmt.Sprintf(netInterfacesNameRow, "iface1") + "\n" +
				fmt.Sprintf(netInterfacesIndexRow, 1) + "\n" +
				netInterfacesAdressessHeader +
				fmt.Sprintf(netInterfacesAdressRow, "ipv4") + "\n" +
				fmt.Sprintf(netInterfacesAdressRow, "ipv6") + "\n" +
				fmt.Sprintf(netInterfacesMTURow, "NA") + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NetInterfacesMessage(tt.netInferfaces)
			if !strings.Contains(got, tt.wantReturnContains) {
				t.Errorf("NetInterfacesMessage() = %v, want contains %v", got, tt.wantReturnContains)
			}
		})
	}
}
