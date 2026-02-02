package handler

import (
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/internal/service"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ByterHandler struct {
	rpc.UnimplementedByterServiceServer
	Service service.ServicerByter
}

func NewByterHandler(srv service.ServicerByter) *ByterHandler {
	return &ByterHandler{Service: srv}
}

func (b *ByterHandler) Upload(stream rpc.ByterService_UploadServer) error {
	obj, err := b.Service.Upload(stream)

	// Какие нибудь спец ошибки ? Internal cлишком //

	if err != nil {
		return status.Errorf(codes.Internal, "%s", err)
	}

	modBytes, ok := obj.(*model.Bytes)
	if !ok {
		return status.Errorf(codes.Internal, "%s", errs.ErrTypeConversion)
	}

	// Response
	err = stream.SendAndClose(&rpc.UploadBytesResponse{
		Status:           "uploaded",
		UpdatedAt:        utils.ConvertDtStr(time.Now()),
		SentFileSize:     modBytes.SentSize,
		ReceivedFileSize: modBytes.ReceivedSize,
	})

	if err != nil {
		return err
	}
	return nil
}

func (b *ByterHandler) Unload(req *rpc.UnloadBytesRequest, stream rpc.ByterService_UnloadServer) error {
	_, err := b.Service.Unload(stream)

	// Запрашиваемые данные не найдены в БД
	if err == errs.ErrDataNotFound {
		return status.Errorf(codes.NotFound, "%s", err)
	}

	// Запрашиваемые данные не найдены в файле
	if err == errs.ErrFileNotFound {
		return status.Errorf(codes.NotFound, "%s", err)
	}

	// Остальные ошибки.
	if err != nil {
		return status.Errorf(codes.Internal, "%s", err)
	}

	return nil
}
