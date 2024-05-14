package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	fmt.Println("Hello Year End Recap!")
	fmt.Println()
	fmt.Println("Here is a random file path:")
	fmt.Println()
  s := filepath.Join("src", "components", "features")
  fmt.Println(s)
}