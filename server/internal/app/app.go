package app

import (
	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/cmd/server"
	"github.com/boginskiy/GophKeeper/server/internal/auth"
	"github.com/boginskiy/GophKeeper/server/internal/db"
	"github.com/boginskiy/GophKeeper/server/internal/handler"
	"github.com/boginskiy/GophKeeper/server/internal/infra"
	"github.com/boginskiy/GophKeeper/server/internal/intercept"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/repo"
	"github.com/boginskiy/GophKeeper/server/internal/service"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
	"github.com/boginskiy/GophKeeper/server/pkg"
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
	fileHandler := utils.NewFileProc()
	pathHandler := utils.NewPathProc()

	// Infra services
	fileManager := infra.NewFileManage(fileHandler, pathHandler)
	cryptoService := pkg.NewCryptoService()

	// Repository
	database := db.NewStoreDB(a.Cfg, a.Logg)
	repoUser := repo.NewRepoUser(a.Cfg, a.Logg, database)
	repoText := repo.NewRepoText(a.Cfg, a.Logg, database)
	repoBytes := repo.NewRepoBytes(a.Cfg, a.Logg, database)

	// Auth
	jwtService := auth.NewJWTService(a.Cfg)
	authService := auth.NewAuth(a.Cfg, a.Logg, jwtService, repoUser)

	// Services
	uploadService := service.NewUploadService(a.Cfg, a.Logg, fileHandler, fileManager, cryptoService, repoBytes)
	bytesService := service.NewBytesService(a.Cfg, a.Logg, repoBytes, fileHandler, fileManager)
	unloadService := service.NewUnloadService(a.Cfg, a.Logg, fileHandler, repoBytes)
	textService := service.NewTextService(a.Cfg, a.Logg, repoText)

	// Interceptor
	interceptor := intercept.NewServIntercept(a.Cfg, a.Logg, authService)

	// Handler
	authHandler := handler.NewAuthHandler(authService)
	texterHandler := handler.NewTexterHandler(textService)
	byterHandler := handler.NewByterHandler(fileHandler, bytesService, unloadService, uploadService)

	// Start server
	server := server.NewServerGRPC(a.Cfg, a.Logg, interceptor)
	server.Registration(authHandler, texterHandler, byterHandler)
	server.Run()

	defer database.CloseDB()
}
