package api

import (
	"context"

	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"google.golang.org/grpc/metadata"
)

type RemoteService struct {
	Cfg       config.Config
	Logg      logg.Logger
	UserChan  chan *model.User
	ClientAPI *ClientAPI
}

func NewRemoteService(
	ctx context.Context,
	config config.Config,
	logger logg.Logger,
	userch chan *model.User,
	clientapi *ClientAPI) *RemoteService {

	tmp := &RemoteService{
		Cfg:       config,
		Logg:      logger,
		UserChan:  userch,
		ClientAPI: clientapi}

	return tmp
}

func (a *RemoteService) Registration(user model.User) (token string, err error) {
	// Request.
	req := &rpc.RegistUserRequest{
		Username:    user.UserName,
		Email:       user.Email,
		Password:    user.Password,
		Phonenumber: user.PhoneNumber}

	// Header.
	var header metadata.MD

	_, err = a.ClientAPI.RegisterUser(req, &header)
	token = a.ClientAPI.TakeValueFromHeader(header, "authorization", 0)

	return token, err
}
