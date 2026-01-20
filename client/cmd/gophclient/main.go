package main

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/app"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

func main() {
	logg := logg.NewLogg("main.log", "INFO")
	cfg := config.NewConf(logg)
	app.NewApp(cfg, logg).Run()
}

// TODO...
// Регать локально можно после получения токена и валидации на сервере
// Add WorkingWithShutdown for server
// Повторить паттерны многопоток
