package helpers

import "slices"

func GetEngineersLeft(allEngineers []string, duplicateEngineers []string) []string {
	final := []string{}
	for _, engineer := range allEngineers {
		if !slices.Contains(duplicateEngineers, engineer) {
			final = append(final, engineer)
		}
	}

	return final
}
