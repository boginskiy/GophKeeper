package app

import (
	"context"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/comm"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/intercept"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
	"github.com/boginskiy/GophKeeper/client/pkg"
)

type App struct {
	Name        string
	Description string
	Version     string
}

func NewApp() *App {
	return &App{
		Name:        config.APPNAME,
		Description: config.DESC,
		Version:     config.VERS,
	}
}

func (a *App) Init(cfg config.Config, logger logg.Logger) {
	// Logger.
	remoteLogger := logg.NewLogService("remote.log", "INFO")

	// Context, channels.
	ctxC, cancel := context.WithCancel(context.Background())
	codeChan := make(chan string, 1)
	mailChan := make(chan string, 1)

	// Utils.
	pathHandler := utils.NewPathProc()
	fileHandler := utils.NewFileProc(pathHandler)
	fileService := infra.NewFileService(fileHandler, pathHandler)

	// User.
	userCLI := user.NewUserCLI(logger)

	// Interceptor.
	interceptor := intercept.NewClientIntercept(cfg, logger, userCLI)

	// Clients & User.
	clientGRPC := client.NewClientGRPC(cfg, logger, interceptor)
	clientCLI := client.NewClientCLI(logger)

	// Infra Services.
	checker := infra.NewCheck(fileHandler)
	dialoger := infra.NewDialog(cfg, logger, checker, clientCLI, userCLI)
	cryptoService := pkg.NewCryptoService()

	// Remote Services.
	remoteAuther := api.NewRemoteAuthService(cfg, remoteLogger, clientGRPC)
	remoteTexter := api.NewRemoteTextService(cfg, remoteLogger, clientGRPC)
	remoteByter := api.NewRemoteBytesService(cfg, remoteLogger, clientGRPC, cryptoService)

	// Business Services.
	bytesService := service.NewBytesService(cfg, logger, fileHandler, pathHandler, remoteByter, fileService)
	textService := service.NewTextService(cfg, logger, remoteTexter)

	// Auth.
	identity := auth.NewIdentity(cfg, logger, fileHandler, pathHandler)
	authService := auth.NewAuthService(cfg, logger, identity, remoteAuther)
	auth.NewRecoveryService(ctxC, cfg, logger, mailChan, codeChan)

	// Commonds.
	commMedia := comm.NewCommMedia(checker, dialoger, bytesService) // Bytes, Sound, Video, Image
	commText := comm.NewCommText(dialoger, textService)

	rootAuth := comm.NewRootAuth(ctxC, mailChan, codeChan, dialoger, authService)
	root := comm.NewRoot(dialoger, commText, commMedia)

	// Start.
	NewRunner(cfg, logger, identity, dialoger, root, rootAuth).Run(clientCLI, userCLI)

	defer clientGRPC.Close()
	defer close(codeChan)
	defer close(mailChan)
	defer logger.Close()
	defer cancel()
}
