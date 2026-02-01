package main

import (
	"fmt"
	"os"
)

func main() {
	currentDir, _ := os.Getwd()

	fmt.Println(currentDir)
}
