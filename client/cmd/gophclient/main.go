package main

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/app"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

func main() {
	mainLogger := logg.NewLogService("main.log", "INFO")
	cfg := config.NewArgsCLI(mainLogger)
	app.NewApp(cfg, mainLogger).Init()
}
