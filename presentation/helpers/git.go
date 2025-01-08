package helpers

import "strings"

func GetReadableCommitMessage(msg string) string {
	return strings.ReplaceAll(msg, "|||", "\n")
}
