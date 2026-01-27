package api

import (
	"context"
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"google.golang.org/grpc/metadata"
)

type RemoteService struct {
	Cfg       config.Config
	Logg      logg.Logger
	ClientAPI *ClientAPI
}

func NewRemoteService(
	ctx context.Context,
	config config.Config,
	logger logg.Logger,
	clientapi *ClientAPI) *RemoteService {

	tmp := &RemoteService{
		Cfg:       config,
		Logg:      logger,
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

func (a *RemoteService) Authentication(user model.User) (token string, err error) {
	req := &rpc.AuthUserRequest{Email: user.Email, Password: user.Password}
	var header metadata.MD

	_, err = a.ClientAPI.AutherUser(req, &header)
	token = a.ClientAPI.TakeValueFromHeader(header, "authorization", 0)
	return token, err
}

func (a *RemoteService) Create(user *user.UserCLI, text model.Text) (any, error) {
	req := &rpc.CreateRequest{
		Name:         text.Name,
		Type:         text.Type,
		Text:         text.Tx,
		Owner:        text.Owner,
		ListActivate: text.ListActivate,
	}

	header := a.ClientAPI.CreateHeaderWithValue("authorization", user.User.Token)
	return a.ClientAPI.Create(req, &header)

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
