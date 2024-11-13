package presentation_helpers

import (
	"encoding/json"
	"fmt"

	"golang.org/x/text/message"
)

func WithCommas(num int) string {
	printer := message.NewPrinter(message.MatchLanguage("en"))

	return printer.Sprintf("%d", num)
}

func TruncateDigits(num float64) string {
	return fmt.Sprintf("%.2f", num)
}

func Json(data any) string {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func Str(data any) string {
	return fmt.Sprintf("%s", data)
}
