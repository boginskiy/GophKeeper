package api

import (
	"context"
	"time"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/manager"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type RemoteAuthService struct {
	Cfg        config.Config
	Logg       logg.Logger
	ClientGRPC *client.ClientGRPC
}

func NewRemoteAuthService(
	config config.Config,
	logger logg.Logger,
	clientgrpc *client.ClientGRPC) *RemoteAuthService {

	return &RemoteAuthService{
		Cfg:        config,
		Logg:       logger,
		ClientGRPC: clientgrpc}
}

func (a *RemoteAuthService) Registration(user model.User) (token string, err error) {
	// Ctx.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.Cfg.GetWaitingTimeResponse()))
	defer cancel()

	// Request.
	req := &rpc.RegistUserRequest{
		Username:    user.UserName,
		Email:       user.Email,
		Password:    user.Password,
		Phonenumber: user.PhoneNumber}

	// Header.
	var serverHeader metadata.MD

	_, err = a.ClientGRPC.AuthService.RegistUser(ctx, req, grpc.Header(&serverHeader))
	token = manager.TakeValueFromHeader(serverHeader, "authorization", 0)

	return token, err
}

func (a *RemoteAuthService) Authentication(user model.User) (token string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.Cfg.GetWaitingTimeResponse()))
	defer cancel()

	req := &rpc.AuthUserRequest{Email: user.Email, Password: user.Password}
	var serverHeader metadata.MD

	_, err = a.ClientGRPC.AuthService.AuthUser(ctx, req, grpc.Header(&serverHeader))
	token = manager.TakeValueFromHeader(serverHeader, "authorization", 0)

	return token, err
}
