package utils

import "fmt"

func PrintStruct(data any) {
	fmt.Println()
	fmt.Printf("%+v\n", data)
	fmt.Println()
}
