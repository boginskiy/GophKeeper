package auth

import (
	"context"

	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/clients"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"google.golang.org/grpc/metadata"
)

type AuthRemote struct {
	Cfg       config.Config
	Logg      logg.Logger
	UserChan  chan *model.User
	ClientAPI *clients.ClientAPI
}

func NewAuthRemote(
	ctx context.Context,
	config config.Config,
	logger logg.Logger,
	userch chan *model.User,
	clientapi *clients.ClientAPI) *AuthRemote {

	tmp := &AuthRemote{
		Cfg:       config,
		Logg:      logger,
		UserChan:  userch,
		ClientAPI: clientapi}

	// Ждем данные для отправки запроса на регистрацию на удаленном сервере.
	go tmp.RemoteRegistration(ctx)

	return tmp
}

func (a *AuthRemote) RemoteRegistration(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case newUser := <-a.UserChan:
			a.UserChan <- a.remoteRegistration(newUser)
		}
	}
}

func (a *AuthRemote) remoteRegistration(user *model.User) *model.User {
	// Request.
	req := &rpc.RegistUserRequest{
		Username:    user.UserName,
		Email:       user.Email,
		Password:    user.Password,
		Phonenumber: user.PhoneNumber}

	// Header.
	var header metadata.MD

	_, err := a.ClientAPI.RegisterUser(req, &header)
	token := a.ClientAPI.TakeValueFromHeader(header, "authorization", 0)

	// Update user.
	user.StatusError = err
	user.Token = token
	return user
}
