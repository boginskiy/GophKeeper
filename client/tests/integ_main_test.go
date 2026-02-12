package tests

// Before tests we must start server $go run ./cmd/gophserver/main.go

import (
	"testing"

	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/stretchr/testify/assert"
)

var (
	dataAuth  = &model.User{Email: "tester@mail.ru", Password: "1234"}
	dataRegis = model.NewUser("Tester", "tester@mail.ru", "89109109910", "1234")
)

func TestClient(t *testing.T) {
	UserCLI := InitUserCLI(logger)
	ServiceAPI := InitServiceAPI(cfg, logger)
	Identy := auth.NewIdentity(cfg, logger, fileHandler)

	ServiceAuth := InitAuthService(cfg, logger, fileHandler, ServiceAPI, Identy)
	ServiceByter := service.NewBytesService(cfg, logger, fileHandler, ServiceAPI)

	Authentication(t, ServiceAuth, UserCLI)

	// Testing
	// testAuthService(t, ServiceAuth, UserCLI)
	testByterService(t, ServiceByter, UserCLI)

	//
	defer Identy.SaveCurrentUser(UserCLI)
}

func Authentication(t *testing.T, service auth.Auth, user user.User) {
	ok := service.Identification(user)
	if ok {
		_, err := service.Authentication(user, dataAuth)
		if err == nil {
			return
		}
	}
	_, err := service.Registration(user, dataRegis)
	assert.NoError(t, err)
}
