package tests

// Before tests we must start server $go run ./cmd/gophserver/main.go

import (
	"context"
	"fmt"
	"testing"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

// TODO!
// Добавить бы автоматическое подниятие сервера, далее прогон тестов
// далее тушим сервер. Так же подумать о тестовой БД

var testLogg = logg.NewLogg("test.log", "INFO")
var testCfg = config.NewConf(
	testLogg, port, attempts, waitTimeRes, retryReq, APPNAME, DESC, VERS, CONFIG)

var testUserCLI = InitUserCLI(testLogg)
var fileHdlr = utils.NewFileHdlr()
var identify = auth.NewIdentity(testCfg, testLogg, fileHdlr)

func TestServiceAPI(t *testing.T) {
	serviceAPI := InitServiceAPI(testCfg, testLogg)

	IsThereRegistr := identify.Identification(testUserCLI)

	if !IsThereRegistr {
		testRegistration(t, serviceAPI)
	} else {
		fmt.Println("FFFF")
	}

	// testAuth(t, serviceAPI)

	identify.SaveCurrentUser(testUserCLI)

}

func InitServiceAPI(cfg *config.Conf, logg *logg.Logg) *api.RemoteService {
	clientGRPC := client.NewClientGRPC(cfg, logg)
	remoteSrv := api.NewRemoteService(context.TODO(), cfg, logg, clientGRPC)
	return remoteSrv
}

func InitUserCLI(logg *logg.Logg) *user.UserCLI {
	userCLI := user.NewUserCLI(logg)
	userCLI.Name = "USER"
	userCLI.User = model.NewUser("Tester", "tester@mail.ru", "89109109910", "1234")
	return userCLI
}
