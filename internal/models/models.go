package models

type CpuInfo struct {
	ID           string
	UsagePercent float64
	CpuSpec
}

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

type HostInfo struct {
	Hostname        string
	OS              string
	Platform        string
	PlatformVersion string
	Uptime          uint64
}
