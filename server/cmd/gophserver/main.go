package main

import (
	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/app"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/pkg"
)

func main() {
	pkg.PrintInfo(config.BuildVersion, config.BuildDate, config.BuildCommit)

	logger := logg.NewLogg("main.log", "INFO")
	cfg := config.NewArgsCLI(logger)
	app.NewApp(cfg, logger).Run()
}
