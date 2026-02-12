package main

import (
	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/app"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/pkg"
)

func main() {
	pkg.PrintInfo(config.BuildVersion, config.BuildDate, config.BuildCommit)

	mainLogger := logg.NewLogService("main.log", "INFO")
	cfg := config.NewArgsCLI(mainLogger)
	app.NewApp(cfg, mainLogger).Run()
}
