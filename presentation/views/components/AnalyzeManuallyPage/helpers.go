package AnalyzeManuallyPage

import "strings"

func GetCombinedValue(nums []string) string {
	strNums := []string{}
	for _, n := range nums {
		strNums = append(strNums, n)
	}

	final := strings.Join(strNums, ",")

	return final
}
