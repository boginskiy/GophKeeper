package app

import (
	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/cmd/server"
	"github.com/boginskiy/GophKeeper/server/internal/auth"
	"github.com/boginskiy/GophKeeper/server/internal/handler"
	"github.com/boginskiy/GophKeeper/server/internal/intercept"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/manager"
	"github.com/boginskiy/GophKeeper/server/internal/repo"
	"github.com/boginskiy/GophKeeper/server/internal/service"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
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
	// Utils.
	fileHandler := utils.NewFileHdlr()

	// Repository
	repoUser := repo.NewRepoUser()
	repoText := repo.NewRepoText()
	repoBytes := repo.NewRepoBytes()

	// Auth
	jwtSrv := auth.NewJWTService(a.Cfg)
	authSrv := auth.NewAuth(a.Cfg, a.Logg, jwtSrv, repoUser)

	// Infra services
	fileManage := manager.NewFileManage(fileHandler)

	// Services
	unloaderSrv := service.NewUnloadService(a.Cfg, a.Logg, fileHandler, repoBytes)

	texterSrv := service.NewTextService(a.Cfg, a.Logg, repoText)
	byterSrv := service.NewBytesService(a.Cfg, a.Logg, repoBytes, fileHandler, fileManage)

	// Interceptor
	interceptor := intercept.NewServIntercept(a.Cfg, a.Logg, authSrv)

	// Handler
	authHdlr := handler.NewAuthHandler(authSrv)
	texterHdlr := handler.NewTexterHandler(texterSrv)
	byterHdlr := handler.NewByterHandler(fileHandler, byterSrv, unloaderSrv)

	// Start server
	server := server.NewServerGRPC(a.Cfg, a.Logg, interceptor)
	server.Registration(authHdlr, texterHdlr, byterHdlr)
	server.Run()
}
