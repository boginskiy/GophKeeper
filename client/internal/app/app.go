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
	Cfg         config.Config
	Logger      logg.Logger
}

func NewApp(conf config.Config, logg logg.Logger) *App {
	tmpApp := &App{
		Cfg:         conf,
		Logger:      logg,
		Name:        config.APPNAME,
		Description: config.DESC,
		Version:     config.VERS,
	}
	return tmpApp
}

func (a *App) Init() {
	// Logger.
	remoteLogger := logg.NewLogService("remote.log", "INFO")

	// ctxT, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctxC, cancel := context.WithCancel(context.Background())
	codeChan := make(chan string, 1)
	mailChan := make(chan string, 1)

	// Utils.
	pathHandler := utils.NewPathProc()
	fileHandler := utils.NewFileProc(pathHandler)
	fileService := infra.NewFileService(fileHandler, pathHandler)

	// User.
	userCLI := user.NewUserCLI(a.Logger)

	// Interceptor.
	interceptor := intercept.NewClientIntercept(a.Cfg, a.Logger, userCLI)

	// Clients & User.
	clientGRPC := client.NewClientGRPC(a.Cfg, a.Logger, interceptor)
	clientCLI := client.NewClientCLI(a.Logger)

	// Infra Services.
	checker := infra.NewCheck(fileHandler)
	dialoger := infra.NewDialog(a.Cfg, a.Logger, checker, clientCLI, userCLI)
	cryptoService := pkg.NewCryptoService()

	// Remote Services.
	remoteAuther := api.NewRemoteAuthService(a.Cfg, remoteLogger, clientGRPC)
	remoteTexter := api.NewRemoteTextService(a.Cfg, remoteLogger, clientGRPC)
	remoteByter := api.NewRemoteBytesService(a.Cfg, remoteLogger, clientGRPC, cryptoService)

	// Business Services.
	bytesService := service.NewBytesService(a.Cfg, a.Logger, fileHandler, pathHandler, remoteByter, fileService)
	textService := service.NewTextService(a.Cfg, a.Logger, remoteTexter)

	// Auth.
	identity := auth.NewIdentity(a.Cfg, a.Logger, fileHandler, pathHandler)
	authService := auth.NewAuthService(a.Cfg, a.Logger, identity, remoteAuther)
	auth.NewRecoveryService(ctxC, a.Cfg, a.Logger, mailChan, codeChan)

	// Commonds.
	commMedia := comm.NewCommMedia(checker, dialoger, bytesService) // Bytes, Sound, Video, Image
	commText := comm.NewCommText(dialoger, textService)

	rootAuth := comm.NewRootAuth(ctxC, mailChan, codeChan, dialoger, authService)
	root := comm.NewRoot(dialoger, commText, commMedia)

	// Start.
	NewRunner(a.Cfg, a.Logger, identity, dialoger, root, rootAuth).Run(clientCLI, userCLI)

	defer clientGRPC.Close()
	defer close(codeChan)
	defer close(mailChan)
	defer a.Logger.Close()
	defer cancel()
}
