package shared

func MapMap[K comparable, A any, B any](mapObj map[K]A, f func(A) B) map[K]B {
	result := make(map[K]B)
	for key, value := range mapObj {
		result[key] = f(value)
	}

	return result
}

func ReduceMap[K comparable, A any](mapObj map[K]A) []A {
	result := []A{}
	for _, value := range mapObj {
		result = append(result, value)
	}

	return result
}
