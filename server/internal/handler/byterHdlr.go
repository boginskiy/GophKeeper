package handler

import (
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/errs"
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

	stream.Context()

	obj, err := b.Srv.Upload(stream)

	if err != nil {
		return status.Errorf(codes.Internal, "%s", err)
	}

	res, ok := obj.(*rpc.UploadBytesResponse)
	if !ok {
		return status.Errorf(codes.Internal, "%s", errs.ErrTypeConversion)
	}

	res.Status = "ok"
	res.UpdatedAt = utils.ConvertDtStr(time.Now())

	// Отправка response to client
	if err := stream.SendAndClose(res); err != nil {
		return err
	}

	return nil
}

// Чтов в итогге на клиенте &
// Пройтись еще раз по цепочке логики
//
