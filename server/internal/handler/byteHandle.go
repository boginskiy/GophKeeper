package handler

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/service"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
)

type ByteHandle struct {
	FileHandler  utils.FileHandler
	BytesService service.BytesServicer[*model.Bytes]
}

func NewByteHandle(
	fileHandler utils.FileHandler,
	bytesService service.BytesServicer[*model.Bytes],
) *ByteHandle {

	return &ByteHandle{
		FileHandler:  fileHandler,
		BytesService: bytesService,
	}
}

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
