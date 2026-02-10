package handler2

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/internal/service"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
)

type ByteHandle struct {
	rpc.UnimplementedByterServiceServer
	FileHandler  utils.FileHandler
	BytesService service.BytesServicer[*model.Bytes]
	// UnloadService service.LoadServicer[rpc.ByterService_UnloadServer, *model.Bytes]
	// UploadService service.LoadServicer[rpc.ByterService_UploadServer, *model.Bytes]
}

func NewByteHandle(
	fileHandler utils.FileHandler,
	bytesService service.BytesServicer[*model.Bytes],
	// unloadService service.LoadServicer[rpc.ByterService_UnloadServer, *model.Bytes],
	// uploadService service.LoadServicer[rpc.ByterService_UploadServer, *model.Bytes],
) *ByteHandle {

	return &ByteHandle{
		FileHandler:  fileHandler,
		BytesService: bytesService,
		// UnloadService: unloadService,
		// UploadService: uploadService
	}
}

// func (b *ByteHandle) Upload(stream rpc.ByterService_UploadServer) error {
// 	modBytes, err := b.UploadService.Prepar(stream)
// 	if err != nil {
// 		return status.Errorf(codes.Internal, "%s", err)
// 	}

// 	modBytes, err = b.UploadService.Load(stream, modBytes)
// 	if err != nil {
// 		return status.Errorf(codes.Internal, "%s", err)
// 	}

// 	// Response
// 	err = stream.SendAndClose(&rpc.UploadBytesResponse{
// 		Status:       "uploaded",
// 		UpdatedAt:    utils.ConvertDtStr(modBytes.UpdatedAt),
// 		SentSize:     modBytes.SentSize,
// 		ReceivedSize: modBytes.ReceivedSize,
// 	})

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (b *ByteHandle) Unload(req *rpc.UnloadBytesRequest, stream rpc.ByterService_UnloadServer) error {
// 	// Получаем данные по файлу из БД.
// 	modBytes, err := b.UnloadService.Prepar(stream)

// 	// Запрашиваемые данные не найдены в БД
// 	if err == errs.ErrDataNotFound {
// 		return status.Errorf(codes.NotFound, "%s", err)
// 	}
// 	// Ошибки связанные с неполучением данных из context.
// 	if err != nil {
// 		return status.Errorf(codes.Internal, "%s", err)
// 	}

// 	// Кладем данные в заголовок для клиента
// 	errUp := infra.PutDataToCtx(stream.Context(), "updated_at", utils.ConvertDtStr(modBytes.UpdatedAt))
// 	errSz := infra.PutDataToCtx(stream.Context(), "sent_size", modBytes.ReceivedSize)
// 	errFn := infra.PutDataToCtx(stream.Context(), "file_name", modBytes.Name)

// 	if errUp != nil || errSz != nil || errFn != nil {
// 		return status.Errorf(codes.Internal, "%s", utils.DefinErr(errUp, errSz, errFn))
// 	}

// 	_, err = b.UnloadService.Load(stream, modBytes)
// 	if err != nil {
// 		return status.Errorf(codes.Internal, "%s", err)
// 	}

// 	return nil
// }

func (b *ByteHandle) HandleRead(ctx context.Context, req *model.Bytes) (*model.Bytes, error) {
	modBytes, err := b.BytesService.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	return modBytes, nil
}

func (b *ByteHandle) HandleReadAll(ctx context.Context, req *model.Bytes) ([]*model.Bytes, error) {
	modBytes, err := b.BytesService.ReadAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return modBytes, nil
}

func (b *ByteHandle) HandleDelete(ctx context.Context, req *model.Bytes) (string, error) {
	_, err := b.BytesService.Delete(ctx, req)
	if err != nil {
		return "", err
	}

	return "deleted", nil
}
