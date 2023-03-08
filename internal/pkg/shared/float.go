package shared

func IsUnequal(a float64, b float64) bool {
	return a < b-0.001 || a > b+0.001
}
