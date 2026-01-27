package api

import (
	"context"
	"fmt"
	"time"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type RemoteService struct {
	Cfg        config.Config
	Logg       logg.Logger
	ClientGRPC *client.ClientGRPC
	// ClientAPI  *ClientAPI
}

func NewRemoteService(
	ctx context.Context,
	config config.Config,
	logger logg.Logger,
	clientgrpc *client.ClientGRPC,
	// clientapi *ClientAPI
) *RemoteService {

	tmp := &RemoteService{
		Cfg:        config,
		Logg:       logger,
		ClientGRPC: clientgrpc,
		// ClientAPI:  clientapi
	}

	return tmp
}

func (a *RemoteService) Registration(user model.User) (token string, err error) {
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
	var header metadata.MD

	_, err = a.ClientGRPC.AuthService.RegistUser(ctx, req, grpc.Header(&header))
	token = a.ClientGRPC.TakeValueFromHeader(header, "authorization", 0)

	return token, err
}

func (a *RemoteService) Authentication(user model.User) (token string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.Cfg.GetWaitingTimeResponse()))
	defer cancel()

	req := &rpc.AuthUserRequest{Email: user.Email, Password: user.Password}
	var header metadata.MD

	_, err = a.ClientGRPC.AuthService.AuthUser(ctx, req, grpc.Header(&header))
	token = a.ClientGRPC.TakeValueFromHeader(header, "authorization", 0)

	return token, err
}

func (a *RemoteService) Create(user *user.UserCLI, text model.Text) (any, error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", user.User.Token)

	req := &rpc.CreateRequest{
		Name:         text.Name,
		Type:         text.Type,
		Text:         text.Tx,
		Owner:        text.Owner,
		ListActivate: text.ListActivate,
	}

	var header metadata.MD

	return a.ClientGRPC.TexterService.Create(ctx, req, grpc.Header(&header))
}

func (a *RemoteService) Read(user *user.UserCLI, text model.Text) {
	// req := "TODO"
	// header := a.ClientAPI.CreateHeaderWithValue("authorization", user.User.Token)
	// req, err := a.ClientAPI.AutherUser(req, &header)
	// что будем возвращать клиенту ?
	fmt.Println("!!! ReadText")
}

func (a *RemoteService) ReadAll(user *user.UserCLI, text model.Text) {
	// req := "TODO"
	// header := a.ClientAPI.CreateHeaderWithValue("authorization", user.User.Token)
	// req, err := a.ClientAPI.AutherUser(req, &header)
	// что будем возвращать клиенту ?
	fmt.Println("!!! ReadText")
}

func (a *RemoteService) Update(user *user.UserCLI, text model.Text) {
	// req := "TODO"
	// header := a.ClientAPI.CreateHeaderWithValue("authorization", user.User.Token)
	// req, err := a.ClientAPI.AutherUser(req, &header)
	// что будем возвращать клиенту ?
	fmt.Println("!!! UpdateText")
}

func (a *RemoteService) Delete(user *user.UserCLI, text model.Text) {
	// req := "TODO"
	// header := a.ClientAPI.CreateHeaderWithValue("authorization", user.User.Token)
	// req, err := a.ClientAPI.AutherUser(req, &header)
	// что будем возвращать клиенту ?
	fmt.Println("!!! Delete")
}
