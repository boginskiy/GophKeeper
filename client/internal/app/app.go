package app

import (
	"bufio"
	"context"
	"os"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/clients"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
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
	mess1Ch := make(chan string, 1)
	mess2Ch := make(chan string, 1)
	UserChan := make(chan *model.User, 1)

	// Logger.
	remoteLogg := logg.NewLogg("remote.log", "INFO")

	// Utils.
	fileHandler := utils.NewWorkingFile()

	// Clients && user
	userCLI := user.NewUserCLI(ctx, a.Logg, mess1Ch, mess2Ch)
	clientCLI := client.NewClientCLI(ctx, mess1Ch, mess2Ch)

	clientGRPC := client.NewClientGRPC(a.Cfg, a.Logg)
	clientAPI := clients.NewClientAPI(a.Cfg, a.Logg, clientGRPC)

	// Auth.
	identity := auth.NewIdentity(a.Logg, fileHandler)
	auther := auth.NewAuth(a.Cfg, a.Logg, UserChan, fileHandler, identity)
	_ = auth.NewAuthServer(ctx, a.Cfg, remoteLogg, UserChan, clientAPI)

	// Services.
	dialogSrv := service.NewDialogService(a.Logg, clientCLI, userCLI, auther)
	dialogSrv.Run(clientCLI, userCLI)

	defer clientGRPC.Close()
	defer a.Logg.Close()
	defer cancel()
}
