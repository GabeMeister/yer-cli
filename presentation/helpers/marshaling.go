package helpers

import (
	"GabeMeister/yer-cli/analyzer"
	"fmt"
	"strings"
)

func MarshalDuplicateGroup(dupGroup analyzer.DuplicateAuthorGroup) string {
	duplicatesStr := strings.Join(dupGroup.Duplicates, "|||")

	// Example: Isaac Neace\tIsaac|||IsaacNeace|||IsaacN
	return fmt.Sprintf("%s\t%s", dupGroup.Real, duplicatesStr)
}

func UnmarshalDuplicateGroup(dupGroupStr string) analyzer.DuplicateAuthorGroup {
	tokens := strings.Split(dupGroupStr, "\t")
	real := tokens[0]
	duplicates := strings.Split(tokens[1], "|||")

	return analyzer.DuplicateAuthorGroup{
		Real:       real,
		Duplicates: duplicates,
	}

}

func MarshalStrSlice(s []string) string {
	return strings.Join(s, ",")
}

func UnmarshalStrSlice(val string) []string {
	values := strings.Split(val, ",")

	if len(values) == 1 && values[0] == "" {
		return []string{}
	}

	return values
}
