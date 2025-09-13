package metrigo

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Matyjash/Metrigo/internal/models"
)

type mockMetricsPuller struct {
	getLogicalCpuCount  func() (int, error)
	getPhysicalCpuCount func() (int, error)
	getCpuUsage         func(bool, time.Duration) ([]float64, error)
	getCpusSpec         func() ([]models.CpuSpec, error)
	getVMMemoryUsage    func() (models.MemoryUsage, error)
	getTemperatures     func() ([]models.TemperatureSensor, error)
}

func (m *mockMetricsPuller) GetLogicalCpuCount() (int, error) {
	return m.getLogicalCpuCount()
}
func (m *mockMetricsPuller) GetPhysicalCpuCount() (int, error) {
	return m.getPhysicalCpuCount()
}
func (m *mockMetricsPuller) GetCpuUsage(percpu bool, interval time.Duration) ([]float64, error) {
	return m.getCpuUsage(percpu, interval)
}
func (m *mockMetricsPuller) GetCpusSpec() ([]models.CpuSpec, error) {
	return m.getCpusSpec()
}
func (m *mockMetricsPuller) GetVMMemoryUsage() (models.MemoryUsage, error) {
	return m.getVMMemoryUsage()
}
func (m *mockMetricsPuller) GetTemperatures() ([]models.TemperatureSensor, error) {
	return m.getTemperatures()
}

// Defaults
var (
	defaultCpuUsage = []float64{10.5, 20.5}
	defaultCpuSpecs = []models.CpuSpec{
		{FrequencyMhz: 3200},
	}
	defaultLogicalCpuCount = 2

	defaultGetLogicalCpuCount = func() (int, error) { return defaultLogicalCpuCount, nil }
	defaultGetCpuUsage        = func(percpu bool, interval time.Duration) ([]float64, error) { return defaultCpuUsage, nil }
	defaultGetCpusSpecs       = func() ([]models.CpuSpec, error) { return defaultCpuSpecs, nil }

	defaultExpectedCpuInfo = []models.CpuInfo{
		{ID: "cpu0", UsagePercent: 10.5, CpuSpec: models.CpuSpec{FrequencyMhz: 3200}},
		{ID: "cpu1", UsagePercent: 20.5, CpuSpec: models.CpuSpec{FrequencyMhz: 3200}},
	}
)

func Test_GetCpuInfo(t *testing.T) {
	tests := []struct {
		name                string
		logicalCpuCountFunc func() (int, error)
		cpuUsageFunc        func(bool, time.Duration) ([]float64, error)
		cpusSpecFunc        func() ([]models.CpuSpec, error)
		wantReturn          []models.CpuInfo
		wantErrContains     string
	}{
		{
			name:       "success with multiple cpus, single cpu spec",
			wantReturn: defaultExpectedCpuInfo,
		},
		{
			name: "success with multiple cpus, multiple cpu specs",
			cpusSpecFunc: func() ([]models.CpuSpec, error) {
				return []models.CpuSpec{
					{FrequencyMhz: 2500},
					{FrequencyMhz: 3200}}, nil
			},
			wantReturn: []models.CpuInfo{
				{ID: "cpu0", UsagePercent: 10.5, CpuSpec: models.CpuSpec{FrequencyMhz: 2500}},
				{ID: "cpu1", UsagePercent: 20.5, CpuSpec: models.CpuSpec{FrequencyMhz: 3200}},
			},
		},
		{
			name:                "success 1 cpu",
			logicalCpuCountFunc: func() (int, error) { return 1, nil },
			cpuUsageFunc:        func(bool, time.Duration) ([]float64, error) { return []float64{99.9}, nil },
			cpusSpecFunc:        func() ([]models.CpuSpec, error) { return []models.CpuSpec{{FrequencyMhz: 2500}}, nil },
			wantReturn: []models.CpuInfo{
				{ID: "cpu0", UsagePercent: 99.9, CpuSpec: models.CpuSpec{FrequencyMhz: 2500}},
			},
		},
		{
			name:                "logical cpu count error",
			logicalCpuCountFunc: func() (int, error) { return 0, fmt.Errorf("fail") },
			wantErrContains:     "failed to get CPU count info",
		},
		{
			name:            "cpu usage error",
			cpuUsageFunc:    func(bool, time.Duration) ([]float64, error) { return nil, fmt.Errorf("fail") },
			wantErrContains: "failed to get CPU usage",
		},
		{
			name:            "cpus info error",
			cpusSpecFunc:    func() ([]models.CpuSpec, error) { return nil, fmt.Errorf("fail") },
			wantErrContains: "failed to get CPUs frequencies",
		},
		{
			name:            "usage length mismatch logical cpu count",
			cpuUsageFunc:    func(bool, time.Duration) ([]float64, error) { return []float64{10}, nil },
			wantErrContains: "mismatched CPU count and usage length",
		},
		{
			// Cpu count: 3, usage count:3, spec count:2
			name:                "spec length not one and not matching logical cpu count",
			logicalCpuCountFunc: func() (int, error) { return 3, nil },
			cpuUsageFunc:        func(bool, time.Duration) ([]float64, error) { return []float64{99.9, 20.51, 47.01}, nil },
			cpusSpecFunc: func() ([]models.CpuSpec, error) {
				return []models.CpuSpec{{FrequencyMhz: 3200}, {FrequencyMhz: 3300}}, nil
			},
			wantErrContains: "not implemented yet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockMetricsPuller{
				getLogicalCpuCount:  defaultGetLogicalCpuCount,
				getCpuUsage:         defaultGetCpuUsage,
				getCpusSpec:         defaultGetCpusSpecs,
				getPhysicalCpuCount: nil,
			}

			if tt.logicalCpuCountFunc != nil {
				mock.getLogicalCpuCount = tt.logicalCpuCountFunc
			}
			if tt.cpuUsageFunc != nil {
				mock.getCpuUsage = tt.cpuUsageFunc
			}
			if tt.cpusSpecFunc != nil {
				mock.getCpusSpec = tt.cpusSpecFunc
			}

			m := Metrigo{}
			m.metricsPuller = mock

			cpus, err := m.GetCpuInfo()
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("expected error containing %q, got \"%v\"", tt.wantErrContains, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tt.wantReturn, cpus) {
				t.Errorf("expected %v, got %v", tt.wantReturn, cpus)
			}
		})
	}
}

func Test_getTemperatures(t *testing.T) {
	tests := []struct {
		name            string
		getTemperatures func() ([]models.TemperatureSensor, error)
		wantReturn      []models.TemperatureSensor
		wantErrContains string
	}{
		{
			name: "successfully gets temperatures sensors",
			getTemperatures: func() ([]models.TemperatureSensor, error) {
				return []models.TemperatureSensor{
					{Key: "sensor1", Value: 45.0},
					{Key: "sensor2", Value: 50.0},
				}, nil
			},
			wantReturn: []models.TemperatureSensor{
				{Key: "sensor1", Value: 45.0},
				{Key: "sensor2", Value: 50.0},
			},
		},
		{
			name: "returns error when getting temperatures fails",
			getTemperatures: func() ([]models.TemperatureSensor, error) {
				return nil, fmt.Errorf("failed to get temperatures")
			},
			wantErrContains: "failed to get temperatures",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockMetricsPuller{
				getTemperatures: tt.getTemperatures,
			}

			m := Metrigo{}
			m.metricsPuller = mock
			temps, err := m.GetTemperatures()
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("expected error containing %q, got \"%v\"", tt.wantErrContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tt.wantReturn, temps) {
				t.Errorf("expected %v, got %v", tt.wantReturn, temps)
			}
		})
	}

}

func Test_GetMemoryUsage(t *testing.T) {
	tests := []struct {
		name             string
		getVMMemoryUsage func() (models.MemoryUsage, error)
		wantReturn       models.MemoryUsage
		wantErrContains  string
	}{
		{
			name: "successfully gets memory usage",
			getVMMemoryUsage: func() (models.MemoryUsage, error) {
				return models.MemoryUsage{UsedB: 1024, TotalB: 2048}, nil
			},
			wantReturn: models.MemoryUsage{UsedB: 1024, TotalB: 2048},
		},
		{
			name: "returns error when getting memory usage fails",
			getVMMemoryUsage: func() (models.MemoryUsage, error) {
				return models.MemoryUsage{}, fmt.Errorf("failed to get memory usage")
			},
			wantErrContains: "failed to get memory usage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockMetricsPuller{
				getVMMemoryUsage: tt.getVMMemoryUsage,
			}
			m := Metrigo{}
			m.metricsPuller = mock
			usage, err := m.GetMemoryUsage()
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("expected error containing %q, got \"%v\"", tt.wantErrContains, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(tt.wantReturn, usage) {
				t.Errorf("expected %v, got %v", tt.wantReturn, usage)
			}
		})
	}
}
