package main

import (
	"log"
	"os"
	"path/filepath"
)

var NAMEFILE = "config.json"
var NAMEAPPL = "gophclient"

func CreateFolderConfig

func ReadOrCreateFileConfig(nameFile string) {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	configPath := filepath.Join(dir, NAMEAPPL, NAMEFILE)
}

func main() {

}
