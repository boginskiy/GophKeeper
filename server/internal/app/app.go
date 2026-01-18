package app

import (
	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/cmd/server"
	"github.com/boginskiy/GophKeeper/server/internal/handler"
	"github.com/boginskiy/GophKeeper/server/internal/intercept"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
)

type App struct {
	Cfg  config.Config
	Logg logg.Logger
}

func NewApp(config config.Config, logger logg.Logger) *App {
	return &App{
		Cfg:  config,
		Logg: logger,
	}
}

func (a *App) Run() {
	// Interceptor
	interceptor := intercept.NewIntercept(a.Cfg, a.Logg)

	// Handler
	keeperHandler := handler.NewKeeperHandler()

	// Start server
	server.NewServerGRPC(a.Cfg, a.Logg).Run(keeperHandler, interceptor)
}
