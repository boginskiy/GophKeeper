package api

import (
	"bufio"
	"context"
	"io"
	"strconv"

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

type RemoteBytesService struct {
	Cfg        config.Config
	Logg       logg.Logger
	ClientGRPC *client.ClientGRPC
}

func NewRemoteBytesService(
	config config.Config,
	logger logg.Logger,
	clientgrpc *client.ClientGRPC) *RemoteBytesService {

	return &RemoteBytesService{
		Cfg:        config,
		Logg:       logger,
		ClientGRPC: clientgrpc}
}

func (a *RemoteBytesService) Upload(user user.User, modBytes model.Bytes) (any, error) {
	var serverHeader metadata.MD

	md := metadata.Pairs(
		"total_size", modBytes.SentSize,
		"file_name", modBytes.Name,
		"file_type", modBytes.Type,
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

// func (a *RemoteBytesService) streamUpload() {

// }

func (a *RemoteBytesService) Unload(user user.User, modBytes model.Bytes) (any, error) {
	md := metadata.Pairs("file_name", modBytes.Name)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var serverHeader metadata.MD

	// Single Request
	stream, err := a.ClientGRPC.ByterService.Unload(ctx, &rpc.UnloadBytesRequest{}, grpc.Header(&serverHeader))
	if err != nil {
		return nil, errs.ErrStartStream.Wrap(err)
	}

	countBytes, err := a.unloadStream(stream, &modBytes)
	if err != nil {
		return nil, err
	}

	// Обогощаем заголовок размером полученных данных
	serverHeader.Set("received_size", strconv.FormatInt(countBytes, 10))

	return &serverHeader, nil
}

func (a *RemoteBytesService) unloadStream(stream rpc.ByterService_UnloadClient, modBytes *model.Bytes) (int64, error) {
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

func (a *RemoteBytesService) Read(user model.User, modBytes model.Bytes) (any, error) {
	md := metadata.Pairs("file_name", modBytes.Name)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var serverHeader metadata.MD

	return a.ClientGRPC.ByterService.Read(ctx, &rpc.ReadBytesRequest{}, grpc.Header(&serverHeader))
}

func (a *RemoteBytesService) ReadAll(user model.User, modBytes model.Bytes) (any, error) {
	md := metadata.Pairs("file_type", modBytes.Type)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var serverHeader metadata.MD

	return a.ClientGRPC.ByterService.ReadAll(ctx, &rpc.ReadAllBytesRequest{}, grpc.Header(&serverHeader))
}

//   rpc Delete (DeleteBytesRequest) returns (DeleteBytesResponse);
// Далее тестирование всего того что сделал!
