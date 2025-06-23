package utils

func MapKeysToSlice(m map[string]bool) []string {
	final := []string{}
	for key := range m {
		final = append(final, key)
	}

	return final
}
