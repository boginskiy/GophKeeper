package app

import (
	"context"

	"github.com/boginskiy/GophKeeper/client/cmd"
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/comm"
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
	// Contexts and channels.
	ctx, cancel := context.WithCancel(context.Background())

	// Logger.
	remoteLogg := logg.NewLogg("remote.log", "INFO")

	// Utils.
	fileHandler := utils.NewWorkingFile()

	// Clients && user
	clientCLI := client.NewClientCLI(a.Logg)
	userCLI := user.NewUserCLI(a.Logg)

	clientGRPC := client.NewClientGRPC(a.Cfg, a.Logg)
	// clientAPI := api.NewClientAPI(a.Cfg, a.Logg, clientGRPC)

	// Infra Services.
	remoteSrv := api.NewRemoteService(ctx, a.Cfg, remoteLogg, clientGRPC)
	dialogSrv := cli.NewDialogService(a.Cfg, a.Logg)

	// Business Services.
	texter := service.NewTexterService(a.Cfg, a.Logg, dialogSrv, remoteSrv)

	// Commonds.
	commImage := comm.NewCommImage(dialogSrv)
	commSound := comm.NewCommSound(dialogSrv)
	commBytes := comm.NewCommBytes(dialogSrv)
	commText := comm.NewCommText(dialogSrv, texter)
	root := comm.NewRoot(dialogSrv, commText, commBytes, commImage, commSound)

	// Auth.
	identity := auth.NewIdentity(a.Logg, fileHandler)
	auther := auth.NewAuth(a.Cfg, a.Logg, fileHandler, identity, dialogSrv, remoteSrv)

	// Start.
	cmd.NewRunner(
		a.Cfg, a.Logg, identity, dialogSrv, auther, root).Run(clientCLI, userCLI)

	defer clientGRPC.Close()
	defer a.Logg.Close()
	defer cancel()
}
