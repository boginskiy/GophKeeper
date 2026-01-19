package app

import (
	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/cmd/server"
	"github.com/boginskiy/GophKeeper/server/internal/auth"
	"github.com/boginskiy/GophKeeper/server/internal/handler"
	"github.com/boginskiy/GophKeeper/server/internal/intercept"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/repository"
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
	// Repository
	repoUser := repository.NewRepoUser()

	// Auth
	jwtSrv := auth.NewJWTService(a.Cfg)
	auther := auth.NewAuth(a.Cfg, a.Logg, jwtSrv, repoUser)

	// Interceptor
	interceptor := intercept.NewIntercept(a.Cfg, a.Logg, auther)

	// Handler
	keeperHandler := handler.NewKeeperHandler(auther)

	// Start server
	server.NewServerGRPC(a.Cfg, a.Logg).Run(keeperHandler, interceptor)
}
