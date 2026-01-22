package app

import (
	"bufio"
	"context"
	"os"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type App struct {
	Name        string
	Description string
	Version     string
	Scanner     *bufio.Scanner
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
		Scanner:     bufio.NewScanner(os.Stdin),
	}

	return tmpApp

}

func (a *App) Run() {
	// Contexts and channels.
	ctx, cancel := context.WithCancel(context.Background())
	messChan := make(chan string, 1)

	// Logger.
	remoteLogg := logg.NewLogg("remote.log", "INFO")

	// Utils.
	fileHandler := utils.NewWorkingFile()

	// Clients && user
	userCLI := user.NewUserCLI(ctx, a.Logg)
	clientCLI := client.NewClientCLI(ctx)

	clientGRPC := client.NewClientGRPC(a.Cfg, a.Logg)
	clientAPI := api.NewClientAPI(a.Cfg, a.Logg, clientGRPC)

	// Services.
	remoteSrv := api.NewRemoteService(ctx, a.Cfg, remoteLogg, clientAPI)
	dialogSrv := cli.NewDialogService(a.Cfg, a.Logg)
	service.NewCommandService(ctx, a.Cfg, a.Logg, messChan)

	// Auth.
	identity := auth.NewIdentity(a.Logg, fileHandler)
	auther := auth.NewAuth(a.Cfg, a.Logg, fileHandler, identity, dialogSrv, remoteSrv)

	service.NewKeeperService(
		a.Cfg, a.Logg, messChan, identity, dialogSrv, auther).Run(clientCLI, userCLI)

	defer clientGRPC.Close()
	defer a.Logg.Close()
	defer cancel()
}
