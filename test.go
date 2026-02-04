package main

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

func main() {

	filePath := "file.txt"

	d := utils.NewFileHdlr()
	fileName := d.FileHandler.TakeFileFromPath(filePath)

	fmt.Println(fileName)
}
