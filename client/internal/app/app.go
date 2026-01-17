package app

import (
	"bufio"
	"context"
	"os"

	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/client"
	"github.com/boginskiy/GophKeeper/client/internal/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/pretty"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type App struct {
	Name        string
	Description string
	Version     string
	Scanner     *bufio.Scanner
	Looker      pretty.Looker
	Logger      logg.Logger
}

func NewApp(logger logg.Logger) *App {
	tmpApp := &App{
		Name:        config.APPNAME,
		Description: config.DESC,
		Version:     config.VERS,
		Scanner:     bufio.NewScanner(os.Stdin),
		Looker:      pretty.NewLook(),
		Logger:      logger,
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

	// Utils.
	fileHandler := utils.NewWorkingFile()

	// Identification.
	identity := user.NewIdentity(a.Logger, fileHandler)
	auth := auth.NewAuth(a.Logger)

	client := client.NewClientCLI(ctx, mess1Ch, mess2Ch)
	user := user.NewUserCLI(ctx, a.Logger, mess1Ch, mess2Ch, identity)

	// Services.
	dialogSrv := service.NewDialogService(a.Logger, client, user, auth)
	dialogSrv.Run(client, user)

	defer user.SaveConfig()
	defer a.Logger.Close()

}
