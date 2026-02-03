package app

import (
	"github.com/boginskiy/GophKeeper/client/cmd"
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/comm"
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
		Name:        conf.GetNameApp(),
		Description: conf.GetDescApp(),
		Version:     conf.GetVersionApp(),
	}
	return tmpApp
}

func (a *App) Init() {
	// Logger.
	remoteLogg := logg.NewLogg("remote.log", "INFO")

	// Utils.
	fileHandler := utils.NewFileHdlr()
	checker := utils.NewCheck()

	// User
	userCLI := user.NewUserCLI(a.Logg)

	// Interceptor
	interceptor := intercept.NewClientIntercept(a.Cfg, a.Logg, userCLI)

	// Clients
	clientGRPC := client.NewClientGRPC(a.Cfg, a.Logg, interceptor)
	clientCLI := client.NewClientCLI(a.Logg)

	// Infra Services.
	dialogSrv := cli.NewDialogService(a.Cfg, a.Logg, checker, clientCLI, userCLI)

	// Remote Services.
	remoteAuthSrv := api.NewRemoteAuthService(a.Cfg, remoteLogg, clientGRPC)
	remoteTextSrv := api.NewRemoteTextService(a.Cfg, remoteLogg, clientGRPC)
	remoteBytesSrv := api.NewRemoteBytesService(a.Cfg, remoteLogg, clientGRPC)

	// Business Services.
	byter := service.NewBytesService(a.Cfg, a.Logg, fileHandler, remoteBytesSrv)
	texter := service.NewTextService(a.Cfg, a.Logg, remoteTextSrv)

	// Auth.
	identity := auth.NewIdentity(a.Cfg, a.Logg, fileHandler)
	authSrv := auth.NewAuthService(a.Cfg, a.Logg, identity, remoteAuthSrv)

	// Commonds.
	commImage := comm.NewCommImage(dialogSrv)
	commSound := comm.NewCommSound(dialogSrv)
	commBytes := comm.NewCommBytes(dialogSrv, byter)
	commText := comm.NewCommText(dialogSrv, texter)
	root := comm.NewRoot(dialogSrv, commText, commBytes, commImage, commSound)

	// Start.
	cmd.NewRunner(
		a.Cfg, a.Logg, identity, dialogSrv, authSrv, root).Run(clientCLI, userCLI)

	defer clientGRPC.Close()
	defer a.Logg.Close()
}
