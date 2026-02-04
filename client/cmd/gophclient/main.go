package main

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/app"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

func main() {
	logg := logg.NewLogg("main.log", "INFO")
	cfg := config.NewArgsCLI(logg)
	app.NewApp(cfg, logg).Init()
}
