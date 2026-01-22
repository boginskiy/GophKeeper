package tests

import (
	"context"
	"testing"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

// Before tests we must start server $go run ./cmd/gophserver/main.go

// TODO!
// Добавить бы автоматическое подниятие сервера, далее прогон тестов
// далее тушим сервер. Так же подумать о тестовой БД

func TestServiceAPI(t *testing.T) {
	logg := logg.NewLogg("test.log", "INFO")
	cfg := config.NewConf(logg)

	serviceAPI := InitServiceAPI(cfg, logg)

	testRegistration(t, serviceAPI)
	testAuthentication(t, serviceAPI)

}

func InitServiceAPI(cfg *config.Conf, logg *logg.Logg) *api.RemoteService {
	clientGRPC := client.NewClientGRPC(cfg, logg)
	clientAPI := api.NewClientAPI(cfg, logg, clientGRPC)
	remoteSrv := api.NewRemoteService(context.TODO(), cfg, logg, clientAPI)

	return remoteSrv
}
