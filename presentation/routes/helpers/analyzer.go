package helpers

import "slices"

func GetAuthorsLeft(allAuthors []string, duplicateAuthors []string) []string {
	final := []string{}
	for _, author := range allAuthors {
		if !slices.Contains(duplicateAuthors, author) {
			final = append(final, author)
		}
	}

	return final
}
