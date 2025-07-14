package utils

func MapKeysToSlice[T comparable](m map[T]bool) []T {
	final := []T{}
	for key := range m {
		final = append(final, key)
	}

	return final
}
