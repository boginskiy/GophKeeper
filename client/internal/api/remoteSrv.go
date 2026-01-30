package api

import (
	"context"
	"io"
	"time"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
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
}

func NewRemoteService(
	ctx context.Context,
	config config.Config,
	logger logg.Logger,
	clientgrpc *client.ClientGRPC) *RemoteService {

	tmp := &RemoteService{
		Cfg:        config,
		Logg:       logger,
		ClientGRPC: clientgrpc}

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

func (a *RemoteService) Create(user user.User, text model.Text) (any, error) {
	ctx := metadata.AppendToOutgoingContext(
		context.Background(), "authorization", user.GetModelUser().Token)

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

func (a *RemoteService) Read(user user.User, text model.Text) (any, error) {
	ctx := metadata.AppendToOutgoingContext(
		context.Background(), "authorization", user.GetModelUser().Token)

	req := &rpc.ReadRequest{Name: text.Name, Owner: text.Owner}
	var header metadata.MD

	return a.ClientGRPC.TexterService.Read(ctx, req, grpc.Header(&header))
}

func (a *RemoteService) ReadAll(user user.User, text model.Text) (any, error) {
	ctx := metadata.AppendToOutgoingContext(
		context.Background(), "authorization", user.GetModelUser().Token)

	req := &rpc.ReadAllRequest{Type: text.Type, Owner: text.Owner}
	var header metadata.MD

	return a.ClientGRPC.TexterService.ReadAll(ctx, req, grpc.Header(&header))
}

func (a *RemoteService) Update(user user.User, text model.Text) (any, error) {
	ctx := metadata.AppendToOutgoingContext(
		context.Background(), "authorization", user.GetModelUser().Token)

	req := &rpc.CreateRequest{
		Name:  text.Name,
		Text:  text.Tx,
		Owner: text.Owner,
	}
	var header metadata.MD
	return a.ClientGRPC.TexterService.Update(ctx, req, grpc.Header(&header))
}

func (a *RemoteService) Delete(user user.User, text model.Text) (any, error) {
	ctx := metadata.AppendToOutgoingContext(
		context.Background(), "authorization", user.GetModelUser().Token)

	req := &rpc.DeleteRequest{
		Name:  text.Name,
		Owner: text.Owner,
	}

	var header metadata.MD
	return a.ClientGRPC.TexterService.Delete(ctx, req, grpc.Header(&header))
}

func (a *RemoteService) Upload(user *user.UserCLI, bytes model.Bytes) (any, error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", user.User.Token)
	var header metadata.MD

	stream, err := a.ClientGRPC.ByterService.Upload(ctx, grpc.Header(&header))
	if err != nil {
		return nil, errs.ErrStartStream.Wrap(err)
	}

	// buffer 1KB.
	buffer := make([]byte, 1024)

	// Run stream.
	for {
		n, err := bytes.Descr.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, errs.ErrReadFileToBuff.Wrap(err)
		}
		if n == 0 {
			break
		}
		// Part of request.
		chunk := &rpc.UploadBytesRequest{
			Content:   buffer[:n],
			TotalSize: bytes.Size,
			Filename:  bytes.Name,
		}
		// Send part of request.
		if err := stream.Send(chunk); err != nil {
			return nil, errs.ErrSendChankFile.Wrap(err)
		}
	}

	return stream.CloseAndRecv()
}
