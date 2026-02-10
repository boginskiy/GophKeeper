package api

import (
	"context"
	"time"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type RemoteAuthService struct {
	Cfg        config.Config
	Logger     logg.Logger
	ClientGRPC *client.ClientGRPC
}

func NewRemoteAuthService(
	config config.Config,
	logger logg.Logger,
	clientgrpc *client.ClientGRPC) *RemoteAuthService {

	return &RemoteAuthService{
		Cfg:        config,
		Logger:     logger,
		ClientGRPC: clientgrpc}
}

func (a *RemoteAuthService) Registration(user model.User) (token string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.Cfg.GetResTimeout()))
	defer cancel()

	// Request.
	req := &rpc.RegistRequest{
		Username:    user.UserName,
		Email:       user.Email,
		Password:    user.Password,
		Phonenumber: user.PhoneNumber}

	// Header.
	var serverHeader metadata.MD

	_, err = a.ClientGRPC.AuthService.Registration(ctx, req, grpc.Header(&serverHeader))
	token = infra.TakeValueFromHeader(serverHeader, "authorization", 0)

	return token, err
}

func (a *RemoteAuthService) Authentication(user model.User) (token string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.Cfg.GetResTimeout()))
	defer cancel()

	req := &rpc.AuthRequest{Email: user.Email, Password: user.Password}
	var serverHeader metadata.MD

	_, err = a.ClientGRPC.AuthService.Authentication(ctx, req, grpc.Header(&serverHeader))
	token = infra.TakeValueFromHeader(serverHeader, "authorization", 0)

	return token, err
}

func (a *RemoteAuthService) Recovery(user model.User) (token string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.Cfg.GetResTimeout()))
	defer cancel()

	// Request.
	req := &rpc.RecovRequest{
		Password: user.Password,
		Email:    user.Email,
	}

	// Header.
	var serverHeader metadata.MD

	_, err = a.ClientGRPC.AuthService.Recovery(ctx, req, grpc.Header(&serverHeader))
	token = infra.TakeValueFromHeader(serverHeader, "authorization", 0)

	return token, err
}
