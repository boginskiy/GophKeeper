package main

import (
	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/app"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
)

func main() {
	logger := logg.NewLogg("LogSrv.txt", "INFO")
	cfg := config.NewConf(logger)
	app.NewApp(cfg, logger).Run()
}
