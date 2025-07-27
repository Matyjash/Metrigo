package metrics

const (
	KiB = 1024
	MiB = 1024 * KiB
	GiB = 1024 * MiB
)

func BytesToGB(bytes uint64) float64 {
	return float64(bytes) / GiB
}
