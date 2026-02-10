package main

import (
	"fmt"
	"time"
)

func main() {

	// Добавляем 3 часа к текущему времени
	timeWithOffset := time.Now().Add(3 * time.Hour).UTC()

	fmt.Println(timeWithOffset)

}
