package main

import (
	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/app"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/pkg"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	pkg.PrintInfo(buildVersion, buildDate, buildCommit)

	logger := logg.NewLogg("main_server.log", "INFO")
	cfg := config.NewConf(logger)
	app.NewApp(cfg, logger).Run()
}
