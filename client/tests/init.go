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
	logger      = logg.NewLogService("test.log", "INFO")
	cfg         = config.NewArgsCLI(logger, port, attempts, waitTimeRes, retryReq, APPNAME, DESC, VERS, CONFIG)
	fileHandler = utils.NewFileProc()
)

func InitServiceAPI(cfg *config.Conf, logger *logg.LogService) *api.RemoteService {
	clientGRPC := client.NewClientGRPC(cfg, logger)
	remoteService := api.NewRemoteService(context.TODO(), cfg, logger, clientGRPC)
	return remoteService
}

func InitUserCLI(logger *logg.LogService) *user.UserCLI {
	userCLI := user.NewUserCLI(logger)
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
