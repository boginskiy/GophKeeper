package handler

import (
	"context"
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/manager"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/internal/service"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ByterHandler struct {
	rpc.UnimplementedByterServiceServer
	FileHdler utils.FileHandler
	BytesSrv  service.BytesServicer[*model.Bytes]
	UnloadSrv service.LoadServicer[rpc.ByterService_UnloadServer, *model.Bytes]
}

func NewByterHandler(
	fileHdler utils.FileHandler,
	bytesSrv service.BytesServicer[*model.Bytes],
	unloadSrv service.LoadServicer[rpc.ByterService_UnloadServer, *model.Bytes]) *ByterHandler {

	return &ByterHandler{FileHdler: fileHdler, BytesSrv: bytesSrv, UnloadSrv: unloadSrv}
}

func (b *ByterHandler) Upload(stream rpc.ByterService_UploadServer) error {
	modBytes, err := b.BytesSrv.Upload(stream)

	if err != nil {
		// Какие нибудь спец ошибки ? Internal cлишком просто //
		return status.Errorf(codes.Internal, "%s", err)
	}

	// Response
	err = stream.SendAndClose(&rpc.UploadBytesResponse{
		Status:       "uploaded",
		UpdatedAt:    utils.ConvertDtStr(modBytes.UpdatedAt),
		SentSize:     modBytes.SentSize,
		ReceivedSize: modBytes.ReceivedSize,
	})

	if err != nil {
		return err
	}
	return nil
}

func (b *ByterHandler) Unload(req *rpc.UnloadBytesRequest, stream rpc.ByterService_UnloadServer) error {
	// Получаем данные по файлу из БД.
	modBytes, err := b.UnloadSrv.Prepar(stream)

	// Запрашиваемые данные не найдены в БД
	if err == errs.ErrDataNotFound {
		return status.Errorf(codes.NotFound, "%s", err)
	}
	// Ошибки связанные с неполучением данных из context.
	if err != nil {
		return status.Errorf(codes.Internal, "%s", err)
	}

	// Кладем данные в заголовок для клиента
	errUp := manager.PutDataToCtx(stream.Context(), "updated_at", utils.ConvertDtStr(modBytes.UpdatedAt))
	errSz := manager.PutDataToCtx(stream.Context(), "sent_size", modBytes.ReceivedSize)
	errFn := manager.PutDataToCtx(stream.Context(), "file_name", modBytes.Name)

	if errUp != nil || errSz != nil || errFn != nil {
		return status.Errorf(codes.Internal, "%s", utils.DefinErr(errUp, errSz, errFn))
	}

	err = b.UnloadSrv.Load(stream, modBytes)
	if err != nil {
		return status.Errorf(codes.Internal, "%s", err)
	}

	return nil
}

func (b *ByterHandler) Read(ctx context.Context, req *rpc.ReadBytesRequest) (*rpc.ReadBytesResponse, error) {
	modBytes, err := b.BytesSrv.Read(ctx, req)

	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// TODO ...
	// Немного костыля :))
	// Сейчас в типе modBytes.Type храниться не тип файла, а идентификатор того, что это данные вида "bytes"
	// Поэтому преобразовываем имя в тип файла, отдаем в response.

	return &rpc.ReadBytesResponse{
		Status:    "read",
		Type:      b.FileHdler.GetTypeFile(modBytes.Name),
		CreatedAt: utils.ConvertDtStr(modBytes.CreatedAt)}, nil
}

func (b *ByterHandler) ReadAll(ctx context.Context, req *rpc.ReadAllBytesRequest) (*rpc.ReadAllBytesResponse, error) {
	modBytes, err := b.BytesSrv.ReadAll(ctx, req)

	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	bytesResponses := make([]*rpc.BytesResponse, len(modBytes))

	for i, bytes := range modBytes {
		bytesResponses[i] = &rpc.BytesResponse{
			UpdatedAt: utils.ConvertDtStr(bytes.UpdatedAt),
			TotalSize: bytes.ReceivedSize,
			Name:      bytes.Name,
		}
	}

	return &rpc.ReadAllBytesResponse{
		Status:         "read",
		BytesResponses: bytesResponses}, nil
}

func (b *ByterHandler) Delete(ctx context.Context, req *rpc.DeleteBytesRequest) (*rpc.DeleteBytesResponse, error) {
	_, err := b.BytesSrv.Delete(ctx, req)

	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	return &rpc.DeleteBytesResponse{
		Status:    "deleted",
		DeletedAt: utils.ConvertDtStr(time.Now())}, nil
}
