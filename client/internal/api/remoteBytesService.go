package api

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type RemoteBytesService struct {
	Cfg           config.Config
	Logger        logg.Logger
	ClientGRPC    *client.ClientGRPC
	CryptoService pkg.Crypter
}

func NewRemoteBytesService(
	config config.Config,
	logger logg.Logger,
	clientgrpc *client.ClientGRPC,
	cryptoService pkg.Crypter) *RemoteBytesService {

	tmp := &RemoteBytesService{
		Cfg:           config,
		Logger:        logger,
		ClientGRPC:    clientgrpc,
		CryptoService: cryptoService}

	// Start CryptoService
	tmp.CryptoService.Start(config.GerCryptoSignature())
	return tmp
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
		fmt.Println(err)
		return nil, errs.ErrStartStream.Wrap(err)
	}

	UploadBytesResponse, err := a.streamUpload(stream, &modBytes)
	return UploadBytesResponse, err
}

func (a *RemoteBytesService) streamUpload(
	stream rpc.ByterService_UploadClient, modBytes *model.Bytes) (*rpc.UploadBytesResponse, error) {
	// Buffer 1KB.
	buffer := make([]byte, 1<<10)
	a.CryptoService.Reset()

	// Run stream.
	for {
		n, err := modBytes.Descr.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, errs.ErrReadFileToBuff.Wrap(err)
		}

		if n == 0 {
			// В конце подписываем данные отправителем для гарантированной доставки файлов.
			chunk := &rpc.UploadBytesRequest{CryptoSignature: a.CryptoService.Sum(nil)}
			if err := stream.Send(chunk); err != nil {
				return nil, pkg.ErrCryptoSignature
			}
			break
		}

		// Part of request.
		chunk := &rpc.UploadBytesRequest{Content: buffer[:n]}
		a.CryptoService.Write(buffer[:n])

		// Send part of request.
		if err := stream.Send(chunk); err != nil {
			return nil, errs.ErrSendChankFile.Wrap(err)
		}
	}
	return stream.CloseAndRecv()
}

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
	// Crypto
	var ClientSignature []byte
	a.CryptoService.Reset()

	var CNT int64

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return CNT, err
		}

		// Take crypto signature from stream.
		if len(req.CryptoSignature) > 0 {
			ClientSignature = req.CryptoSignature
		}

		nn, err := writer.Write(req.Content)
		a.CryptoService.Write(req.Content)

		if err != nil {
			return CNT, err
		}

		CNT += int64(nn)
	}

	// Check signature.
	ok := a.CryptoService.CheckSignature(ClientSignature)
	if !ok {
		return 0, pkg.ErrCheckCryptoSignature
	}

	err := writer.Flush()
	if err != nil {
		return CNT, err
	}

	return CNT, nil
}

func (a *RemoteBytesService) Read(user user.User, modBytes model.Bytes) (any, error) {
	md := metadata.Pairs("file_name", modBytes.Name)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var serverHeader metadata.MD

	return a.ClientGRPC.ByterService.Read(ctx, &rpc.ReadBytesRequest{}, grpc.Header(&serverHeader))
}

func (a *RemoteBytesService) ReadAll(user user.User, modBytes model.Bytes) (any, error) {
	md := metadata.Pairs("file_type", modBytes.Type)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var serverHeader metadata.MD

	return a.ClientGRPC.ByterService.ReadAll(ctx, &rpc.ReadAllBytesRequest{}, grpc.Header(&serverHeader))
}

func (a *RemoteBytesService) Delete(user user.User, modBytes model.Bytes) (any, error) {
	md := metadata.Pairs("file_name", modBytes.Name)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var serverHeader metadata.MD

	return a.ClientGRPC.ByterService.Delete(ctx, &rpc.DeleteBytesRequest{}, grpc.Header(&serverHeader))
}
