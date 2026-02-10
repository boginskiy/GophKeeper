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
	"github.com/boginskiy/GophKeeper/server/internal/handler2"
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
	repoUser := repo.NewRepoUser(a.Cfg, a.Logger, database)
	repoText := repo.NewRepoText(a.Cfg, a.Logger, database)
	repoBytes := repo.NewRepoBytes(a.Cfg, a.Logger, database)

	// Auth
	jwtService := auth.NewJWTService(a.Cfg)
	authService := auth.NewAuth(a.Cfg, a.Logger, jwtService, repoUser)

	// Services
	uploadService := service.NewUploadService(a.Cfg, a.Logger, fileHandler, fileService, cryptoService, repoBytes)
	bytesService := service.NewBytesService(a.Cfg, a.Logger, repoBytes, fileHandler, fileService)
	unloadService := service.NewUnloadService(a.Cfg, a.Logger, fileHandler, repoBytes)
	textService := service.NewTextService(a.Cfg, a.Logger, repoText)

	// Interceptor
	interceptor := intercept.NewServIntercept(a.Cfg, a.Logger, authService)

	// Codec
	byteCoder := codec.NewByteDecoderEncoder(fileService, fileHandler, repoBytes)
	authDecoder := codec.NewAuthDecoderEncoder()

	// Handler
	// authHandler := handler.NewAuthHandler(authService)
	texterHandler := handler.NewTexterHandler(textService)
	// byterHandler := handler.NewByterHandler(fileHandler, bytesService, unloadService, uploadService)

	// Handler 2
	authHandler := handler2.NewAuthHandle(authService)
	// textHandler := handler2.NewTexterHandler(textService)
	byteHandler := handler2.NewByteHandle(fileHandler, bytesService)

	// Facades
	authFacade := facade.NewAuthFacade(errMapper, authHandler, authDecoder)
	byteFacade := facade.NewByteFacade(errMapper, byteCoder, byteHandler, uploadService, unloadService)

	// Start server
	server := server.NewServerGRPC(a.Cfg, a.Logger, interceptor)
	server.Registration(authFacade, texterHandler, byteFacade)
	server.Run()

	defer database.CloseDB()
}

// ErrMapper     errs.ErrMapper
// 	ByteCoder     codec.ByteGRPCCoder[*model.Bytes]
// 	UploadService service.LoadServicer[rpc.ByterService_UploadServer, *model.Bytes]
