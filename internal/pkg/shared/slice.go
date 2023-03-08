package shared

func MapSlice[A any, B any](slice []A, f func(A) B) []B {
	result := make([]B, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}

	return result
}
