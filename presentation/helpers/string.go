package presentation_helpers

import (
	"encoding/json"

	"golang.org/x/text/message"
)

func WithCommas(num int) string {
	printer := message.NewPrinter(message.MatchLanguage("en"))

	return printer.Sprintf("%d", num)
}

func Json(data any) string {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bytes)
}