package api

import (
	"bufio"
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
	var serverHeader metadata.MD

	_, err = a.ClientGRPC.AuthService.RegistUser(ctx, req, grpc.Header(&serverHeader))
	token = a.ClientGRPC.TakeValueFromHeader(serverHeader, "authorization", 0)

	return token, err
}

func (a *RemoteService) Authentication(user model.User) (token string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.Cfg.GetWaitingTimeResponse()))
	defer cancel()

	req := &rpc.AuthUserRequest{Email: user.Email, Password: user.Password}
	var serverHeader metadata.MD

	_, err = a.ClientGRPC.AuthService.AuthUser(ctx, req, grpc.Header(&serverHeader))
	token = a.ClientGRPC.TakeValueFromHeader(serverHeader, "authorization", 0)

	return token, err
}

func (a *RemoteService) Create(user user.User, text model.Text) (any, error) {
	var serverHeader metadata.MD

	req := &rpc.CreateRequest{
		Name:         text.Name,
		Type:         text.Type,
		Text:         text.Tx,
		Owner:        text.Owner,
		ListActivate: text.ListActivate,
	}

	return a.ClientGRPC.TexterService.Create(context.Background(), req, grpc.Header(&serverHeader))
}

func (a *RemoteService) Read(user user.User, text model.Text) (any, error) {
	var serverHeader metadata.MD
	req := &rpc.ReadRequest{Name: text.Name, Owner: text.Owner}

	return a.ClientGRPC.TexterService.Read(context.Background(), req, grpc.Header(&serverHeader))
}

func (a *RemoteService) ReadAll(user user.User, text model.Text) (any, error) {
	var serverHeader metadata.MD
	req := &rpc.ReadAllRequest{Type: text.Type, Owner: text.Owner}

	return a.ClientGRPC.TexterService.ReadAll(context.Background(), req, grpc.Header(&serverHeader))
}

func (a *RemoteService) Update(user user.User, text model.Text) (any, error) {
	var serverHeader metadata.MD

	req := &rpc.CreateRequest{
		Name:  text.Name,
		Text:  text.Tx,
		Owner: text.Owner,
	}

	return a.ClientGRPC.TexterService.Update(context.Background(), req, grpc.Header(&serverHeader))
}

func (a *RemoteService) Delete(user user.User, text model.Text) (any, error) {
	var serverHeader metadata.MD

	req := &rpc.DeleteRequest{
		Name:  text.Name,
		Owner: text.Owner,
	}
	return a.ClientGRPC.TexterService.Delete(context.Background(), req, grpc.Header(&serverHeader))
}

func (a *RemoteService) Upload(user user.User, modBytes model.Bytes) (any, error) {
	var serverHeader metadata.MD

	md := metadata.Pairs(
		"total_size", modBytes.SentSize,
		"file_name", modBytes.Name,
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Stream.
	stream, err := a.ClientGRPC.ByterService.Upload(ctx, grpc.Header(&serverHeader))
	if err != nil {
		return nil, errs.ErrStartStream.Wrap(err)
	}

	// Buffer 1KB.
	buffer := make([]byte, 1024)

	// Run stream.
	for {
		n, err := modBytes.Descr.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, errs.ErrReadFileToBuff.Wrap(err)
		}
		if n == 0 {
			break
		}
		// Part of request.
		chunk := &rpc.UploadBytesRequest{Content: buffer[:n]}

		// Send part of request.
		if err := stream.Send(chunk); err != nil {
			return nil, errs.ErrSendChankFile.Wrap(err)
		}
	}
	return stream.CloseAndRecv()
}

// func (a *RemoteService) streamUpload() {

// }

func (a *RemoteService) Unload(user user.User, modBytes model.Bytes) (any, error) {
	md := metadata.Pairs("file_name", modBytes.Name)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var serverHeader metadata.MD

	// Single Request
	stream, err := a.ClientGRPC.ByterService.Unload(ctx, &rpc.UnloadBytesRequest{}, grpc.Header(&serverHeader))
	if err != nil {
		return nil, errs.ErrStartStream.Wrap(err)
	}

	// TODO имя, размер, обновление

	_, err = a.unloadStream(stream, &modBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *RemoteService) unloadStream(stream rpc.ByterService_UnloadClient, modBytes *model.Bytes) (int64, error) {
	// Writer
	writer := bufio.NewWriter(modBytes.Descr)
	var CNT int64

	for {
		// Обработка запроса.
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return CNT, err
		}

		nn, err := writer.Write(req.Content)
		if err != nil {
			return CNT, err
		}

		CNT += int64(nn)
	}

	err := writer.Flush()
	if err != nil {
		return CNT, err
	}

	return CNT, nil
}
