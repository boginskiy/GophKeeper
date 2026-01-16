package app

import (
	"bufio"
	"context"
	"os"

	"github.com/boginskiy/GophKeeper/client/internal/client"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/pretty"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

const (
	NAME = "gophclient"
	DESC = "HI man! It is special CLI application for your computer"
	VERS = "1.1.01"
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
		Name:        NAME,
		Description: DESC,
		Version:     VERS,
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

	// Utils
	fileHandler := utils.NewWorkingFile()

	// System identification for user.
	identity := user.NewIdentity(a.Logger, fileHandler)

	client := client.NewClientCLI(ctx, mess1Ch, mess2Ch)
	user := user.NewUserCLI(ctx, a.Logger, mess1Ch, mess2Ch, identity)

	// Services
	dialogSrv := service.NewDialogService(a.Logger, client, user)

	dialogSrv.Run()

	defer user.SaveConfig()

}

// TODO...
// Используйте bufio.Scanner для простого ввода одной строки.
// Применяйте bufio.Reader для чтения нескольких строк подряд.
// Воспользуйтесь fmt.Scanf для ввода конкретных типов данных (число, строка и т.д.).
// Используйте пакет term для безопасной передачи конфиденциальных данных (паролей).
