package app

import (
	"context"
	"time"

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
)

type App struct {
	Name        string
	Description string
	Version     string
	Cfg         config.Config
	Logg        logg.Logger
}

func NewApp(conf config.Config, logg logg.Logger) *App {
	tmpApp := &App{
		Cfg:         conf,
		Logg:        logg,
		Name:        config.APPNAME,
		Description: config.DESC,
		Version:     config.VERS,
	}
	return tmpApp
}

func (a *App) Init() {
	// Logger.
	remoteLogg := logg.NewLogg("remote.log", "INFO")

	ctxT, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctxC, cancel := context.WithCancel(context.Background())
	codeChan := make(chan string, 1)
	mailChan := make(chan string, 1)

	// Utils.
	pathHandler := utils.NewPathProc()
	fileHandler := utils.NewFileProc(pathHandler)
	fileManager := infra.NewFileManage(fileHandler, pathHandler)

	// User.
	userCLI := user.NewUserCLI(a.Logg)

	// Interceptor.
	interceptor := intercept.NewClientIntercept(a.Cfg, a.Logg, userCLI)

	// Clients & User.
	clientGRPC := client.NewClientGRPC(a.Cfg, a.Logg, interceptor)
	clientCLI := client.NewClientCLI(a.Logg)

	// Infra Services.
	checker := infra.NewCheck(fileHandler)
	dialoger := infra.NewDialog(a.Cfg, a.Logg, checker, clientCLI, userCLI)

	// Remote Services.
	remoteAuther := api.NewRemoteAuthService(a.Cfg, remoteLogg, clientGRPC)
	remoteTexter := api.NewRemoteTextService(a.Cfg, remoteLogg, clientGRPC)
	remoteByter := api.NewRemoteBytesService(a.Cfg, remoteLogg, clientGRPC)

	// Business Services.
	bytesService := service.NewBytesService(a.Cfg, a.Logg, fileHandler, pathHandler, remoteByter, fileManager)
	textService := service.NewTextService(a.Cfg, a.Logg, remoteTexter)

	// Auth.
	identity := auth.NewIdentity(a.Cfg, a.Logg, fileHandler, pathHandler)
	authSrv := auth.NewAuthService(a.Cfg, a.Logg, identity, remoteAuther)
	auth.NewRecovery(ctxC, a.Cfg, a.Logg, mailChan, codeChan)

	// Commonds.
	commMedia := comm.NewCommMedia(checker, dialoger, bytesService) // Bytes, Sound, Video, Image
	commText := comm.NewCommText(dialoger, textService)
	root := comm.NewRoot(ctxT, dialoger, commText, commMedia, mailChan, codeChan)

	// Start.
	NewRunner(
		a.Cfg, a.Logg, identity, dialoger, authSrv, root).Run(clientCLI, userCLI)

	defer clientGRPC.Close()
	defer close(codeChan)
	defer close(mailChan)
	defer a.Logg.Close()
	defer cancel()
}
