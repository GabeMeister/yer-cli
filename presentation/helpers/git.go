package presentation_helpers

import "strings"

func GetReadableCommitMessage(msg string) string {
	return strings.ReplaceAll(msg, "|||", "\n")
}
