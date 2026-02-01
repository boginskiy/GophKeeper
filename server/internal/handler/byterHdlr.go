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
	Srv service.ServicerByter
}

func NewByterHandler(srv service.ServicerByter) *ByterHandler {
	return &ByterHandler{Srv: srv}
}

func (b *ByterHandler) Upload(stream rpc.ByterService_UploadServer) error {
	obj, err := b.Srv.Upload(stream)

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
