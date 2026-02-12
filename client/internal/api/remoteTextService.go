package api

import (
	"context"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type RemoteTextService struct {
	Cfg        config.Config
	Logger     logg.Logger
	ClientGRPC *client.ClientGRPC
}

func NewRemoteTextService(
	config config.Config,
	logger logg.Logger,
	clientgrpc *client.ClientGRPC) *RemoteTextService {

	return &RemoteTextService{
		Cfg:        config,
		Logger:     logger,
		ClientGRPC: clientgrpc}
}

func (a *RemoteTextService) Create(user user.User, text model.Text) (any, error) {
	var serverHeader metadata.MD

	req := &rpc.CreateRequest{
		Name:         text.Name,
		Type:         text.Type,
		Text:         text.Content,
		Owner:        text.Owner,
		ListActivate: text.ListActivate,
	}

	return a.ClientGRPC.TexterService.Create(context.Background(), req, grpc.Header(&serverHeader))
}

func (a *RemoteTextService) Read(user user.User, text model.Text) (any, error) {
	var serverHeader metadata.MD
	req := &rpc.ReadRequest{Name: text.Name, Owner: text.Owner}

	return a.ClientGRPC.TexterService.Read(context.Background(), req, grpc.Header(&serverHeader))
}

func (a *RemoteTextService) ReadAll(user user.User, text model.Text) (any, error) {
	var serverHeader metadata.MD
	req := &rpc.ReadAllRequest{Type: text.Type, Owner: text.Owner}

	return a.ClientGRPC.TexterService.ReadAll(context.Background(), req, grpc.Header(&serverHeader))
}

func (a *RemoteTextService) Update(user user.User, text model.Text) (any, error) {
	var serverHeader metadata.MD

	req := &rpc.CreateRequest{
		Name:  text.Name,
		Text:  text.Content,
		Owner: text.Owner,
	}

	return a.ClientGRPC.TexterService.Update(context.Background(), req, grpc.Header(&serverHeader))
}

func (a *RemoteTextService) Delete(user user.User, text model.Text) (any, error) {
	var serverHeader metadata.MD

	req := &rpc.DeleteRequest{
		Name:  text.Name,
		Owner: text.Owner,
	}
	return a.ClientGRPC.TexterService.Delete(context.Background(), req, grpc.Header(&serverHeader))
}
