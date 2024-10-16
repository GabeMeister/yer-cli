package utils

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func FindIndex[T any](slice []T, test func(T) bool) int {
	for idx, s := range slice {
		if test(s) {
			return idx
		}
	}

	return -1
}

func Includes[T any](slice []T, test func(T) bool) bool {
	for _, s := range slice {
		if test(s) {
			return true
		}
	}

	return false
}

func DeleteAtIndex(slice []any, index int) []any {
	return append(slice[:index], slice[index+1:])
}

func Delete[T any](slice []T, test func(T) bool) (final []T) {
	for _, item := range slice {
		if !test(item) {
			final = append(final, item)
		}
	}

	return final
}
