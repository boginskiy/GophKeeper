package app

import (
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

	// Utils.
	fileHandler := utils.NewFileHdlr()

	// User.
	userCLI := user.NewUserCLI(a.Logg)

	// Interceptor.
	interceptor := intercept.NewClientIntercept(a.Cfg, a.Logg, userCLI)

	// Clients.
	clientGRPC := client.NewClientGRPC(a.Cfg, a.Logg, interceptor)
	clientCLI := client.NewClientCLI(a.Logg)

	// Infra Services.
	checker := infra.NewCheck(fileHandler)
	dialogSrv := infra.NewDialogService(a.Cfg, a.Logg, checker, clientCLI, userCLI)

	// Remote Services.
	remoteAuthSrv := api.NewRemoteAuthService(a.Cfg, remoteLogg, clientGRPC)
	remoteTextSrv := api.NewRemoteTextService(a.Cfg, remoteLogg, clientGRPC)
	remoteBytesSrv := api.NewRemoteBytesService(a.Cfg, remoteLogg, clientGRPC)

	// Business Services.
	byterSrv := service.NewBytesService(a.Cfg, a.Logg, fileHandler, remoteBytesSrv)
	texterSrv := service.NewTextService(a.Cfg, a.Logg, remoteTextSrv)

	// Auth.
	identity := auth.NewIdentity(a.Cfg, a.Logg, fileHandler)
	authSrv := auth.NewAuthService(a.Cfg, a.Logg, identity, remoteAuthSrv)

	// Commonds.
	commMedia := comm.NewCommMedia(checker, dialogSrv, byterSrv) // Bytes, Sound, Video, Image
	commText := comm.NewCommText(dialogSrv, texterSrv)

	root := comm.NewRoot(dialogSrv, commText, commMedia)

	// Start.
	NewRunner(
		a.Cfg, a.Logg, identity, dialogSrv, authSrv, root).Run(clientCLI, userCLI)

	defer clientGRPC.Close()
	defer a.Logg.Close()
}
