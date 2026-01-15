package app

import (
	"bufio"
	"context"
	"os"

	"github.com/boginskiy/GophKeeper/client/internal/client"
	"github.com/boginskiy/GophKeeper/client/internal/pretty"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

const (
	NAME = "GophClient"
	DESC = "HI man! It is special CLI application for your computer"
	VERS = "1.1.01"
)

type App struct {
	Name        string
	Description string
	Version     string
	Scanner     *bufio.Scanner
	Looker      pretty.Looker
}

func NewApp() *App {
	tmpApp := &App{
		Name:        NAME,
		Description: DESC,
		Version:     VERS,
		Scanner:     bufio.NewScanner(os.Stdin),
		Looker:      pretty.NewLook(),
	}

	// tmpApp.Looker.Hello(pretty.ClientCLI, pretty.UserCLI)

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

	client := client.NewClientCLI(ctx, "CLIENT", mess1Ch, mess2Ch)
	user := user.NewUserCLI(ctx, "USER", mess1Ch, mess2Ch)

	// Services
	dialogSrv := service.NewDialogService(client, user)

	dialogSrv.Run()

	// Приветствие
	// a.Looker.Hello(pretty.ClientCLI, pretty.UserCLI)

	// for {
	// 	a.Looker.PrintInfo(pretty.ClientCLI, `Enter the command or 'help'...`)
	// 	a.Scanner.Scan()
	// }
}

// TODO...
// Используйте bufio.Scanner для простого ввода одной строки.
// Применяйте bufio.Reader для чтения нескольких строк подряд.
// Воспользуйтесь fmt.Scanf для ввода конкретных типов данных (число, строка и т.д.).
// Используйте пакет term для безопасной передачи конфиденциальных данных (паролей).
