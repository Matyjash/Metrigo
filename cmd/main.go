package main

import (
	"github.com/Matyjash/Metrigo/internal/metrigo"
)

func main() {
	metrigo := metrigo.NewMetrigo()
	_, _ = metrigo.GetCpuInfo()
}
