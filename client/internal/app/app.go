package app

import (
	"bufio"
	"context"
	"os"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
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
	// Config
	// Logger
	//

	mess1Ch := make(chan string, 1)
	mess2Ch := make(chan string, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// gRPC
	clientGRPC := client.NewClientGRPC(a.Cfg, a.Logg).Run()

	// Utils.
	fileHandler := utils.NewWorkingFile()

	// Identification.
	identity := user.NewIdentity(a.Logg, fileHandler)
	auth := auth.NewAuth(a.Logg)

	clientCLI := client.NewClientCLI(ctx, mess1Ch, mess2Ch)
	userCLI := user.NewUserCLI(ctx, a.Logg, mess1Ch, mess2Ch, identity)

	// Services.
	dialogSrv := service.NewDialogService(a.Logg, clientCLI, userCLI, auth)
	dialogSrv.Run(clientCLI, userCLI)

	defer userCLI.SaveConfig()
	defer a.Logg.Close()

}
