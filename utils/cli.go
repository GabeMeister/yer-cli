package utils

import (
	"bufio"
	"fmt"
	"os"

	"github.com/briandowns/spinner"
)

func Pause(vals ...interface{}) {
	fmt.Println()
	if len(vals) > 0 {
		fmt.Println(vals...)
	}
	fmt.Println()

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func IsDevMode() bool {
	return os.Getenv("DEV_MODE") == "true"
}

func PrintProgress(s *spinner.Spinner, msg string) {
	if IsDevMode() {
		fmt.Println(msg)
	} else {
		s.Suffix = " " + msg
	}
}
