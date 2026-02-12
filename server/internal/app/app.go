package app

import (
	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/cmd/server"
	"github.com/boginskiy/GophKeeper/server/internal/auth"
	"github.com/boginskiy/GophKeeper/server/internal/codec"
	"github.com/boginskiy/GophKeeper/server/internal/db"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/facade"
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
	Cfg    config.Config
	Logger logg.Logger
}

func NewApp(config config.Config, logger logg.Logger) *App {
	return &App{
		Cfg:    config,
		Logger: logger,
	}
}

func (a *App) Run() {
	// Utils.
	fileHandler := utils.NewFileProc()
	pathHandler := utils.NewPathProc()

	// Infra services
	fileService := infra.NewFileService(fileHandler, pathHandler)
	cryptoService := pkg.NewCryptoService()
	errMapper := errs.NewErrMap()

	// Repository
	database := db.NewStoreDB(a.Cfg, a.Logger)
	repoBytes := repo.NewRepoBytes(a.Cfg, a.Logger, database)
	repoUser := repo.NewRepoUser(a.Cfg, a.Logger, database)
	repoText := repo.NewRepoText(a.Cfg, a.Logger, database)

	// Auth
	jwtService := auth.NewJWTService(a.Cfg)
	authService := auth.NewAuth(a.Cfg, a.Logger, jwtService, repoUser)

	// Services
	uploadService := service.NewUploadService(a.Cfg, a.Logger, fileHandler, fileService, cryptoService, repoBytes)
	unloadService := service.NewUnloadService(a.Cfg, a.Logger, fileHandler, cryptoService, repoBytes)
	bytesService := service.NewBytesService(a.Cfg, a.Logger, repoBytes, fileHandler, fileService)
	textService := service.NewTextService(a.Cfg, a.Logger, repoText)

	// Interceptor
	interceptor := intercept.NewServIntercept(a.Cfg, a.Logger, authService)

	// Codec
	byteCoder := codec.NewByteDecoderEncoder(fileService, fileHandler, repoBytes)
	authDecoder := codec.NewAuthDecoderEncoder()
	textCoder := codec.NewTextDecoderEncoder()

	// Handler
	byteHandler := handler.NewByteHandle(fileHandler, bytesService)
	authHandler := handler.NewAuthHandle(authService)
	textHandler := handler.NewTextHandle(textService)

	// Facades
	byteFacade := facade.NewByteFacade(errMapper, byteCoder, byteHandler, uploadService, unloadService)
	authFacade := facade.NewAuthFacade(errMapper, authHandler, authDecoder)
	textFacade := facade.NewTextFacade(errMapper, textCoder, textHandler)

	// Start server
	server := server.NewServerGRPC(a.Cfg, a.Logger, interceptor)
	server.Registration(authFacade, textFacade, byteFacade)
	server.Run()

	defer database.CloseDB()
}
