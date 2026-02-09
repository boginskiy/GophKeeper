package tests

import (
	"context"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

var (
	logger      = logg.NewLogg("test.log", "INFO")
	cfg         = config.NewArgsCLI(logger, port, attempts, waitTimeRes, retryReq, APPNAME, DESC, VERS, CONFIG)
	fileHandler = utils.NewFileProc()
)

func InitServiceAPI(cfg *config.Conf, logg *logg.Logg) *api.RemoteService {
	clientGRPC := client.NewClientGRPC(cfg, logg)
	remoteService := api.NewRemoteService(context.TODO(), cfg, logg, clientGRPC)
	return remoteService
}

func InitUserCLI(logg *logg.Logg) *user.UserCLI {
	userCLI := user.NewUserCLI(logg)
	userCLI.Name = "USER"
	userCLI.User = model.NewUser("Tester", "tester@mail.ru", "89109109910", "1234")
	return userCLI
}

func InitAuthService(
	cfg config.Config,
	logger logg.Logger,
	fileHandler utils.FileHandler,
	api api.RemoteAuther,
	identy auth.Identifier,

) *auth.AuthService {
	return auth.NewAuthService(cfg, logger, fileHandler, identy, api)
}
