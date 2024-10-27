package utils

import (
	"bufio"
	"fmt"
	"os"
)

func Pause(vals ...interface{}) {
	if len(vals) > 0 {
		fmt.Println(vals...)
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
