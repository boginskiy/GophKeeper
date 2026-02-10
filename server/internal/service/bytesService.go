package service

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/infra"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/repo"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
)

type BytesService struct {
	Cfg         config.Config
	Logger      logg.Logger
	Repo        repo.Repository[*model.Bytes]
	FileHandler utils.FileHandler
	FileService infra.Filer
}

func NewBytesService(
	config config.Config,
	logger logg.Logger,
	repo repo.Repository[*model.Bytes],
	fileHandler utils.FileHandler,
	fileService infra.Filer) *BytesService {

	tmp := &BytesService{
		Cfg:         config,
		Logger:      logger,
		Repo:        repo,
		FileHandler: fileHandler,
		FileService: fileService,
	}

	return tmp
}

func (b *BytesService) Read(ctx context.Context, _ any) (*model.Bytes, error) {
	// Info from context.
	fileName, err := infra.TakeClientValueFromCtx(ctx, "file_name", 0)
	if err != nil {
		return nil, err
	}
	owner, err := infra.TakeServerValStrFromCtx(ctx, infra.EmailCtx)
	if err != nil {
		return nil, err
	}

	return b.Repo.ReadRecord(ctx, &model.Bytes{Name: fileName, Owner: owner})
}

func (b *BytesService) ReadAll(ctx context.Context, req any) ([]*model.Bytes, error) {
	// Info from context.
	fileType, err := infra.TakeClientValueFromCtx(ctx, "file_type", 0)
	if err != nil {
		return nil, err
	}
	owner, err := infra.TakeServerValStrFromCtx(ctx, infra.EmailCtx)
	if err != nil {
		return nil, err
	}

	return b.Repo.ReadAllRecord(ctx, &model.Bytes{Type: fileType, Owner: owner})
}

func (b *BytesService) Delete(ctx context.Context, req any) (*model.Bytes, error) {
	// Info from context.
	fileName, err := infra.TakeClientValueFromCtx(ctx, "file_name", 0)
	if err != nil {
		return nil, err
	}
	owner, err := infra.TakeServerValStrFromCtx(ctx, infra.EmailCtx)
	if err != nil {
		return nil, err
	}

	return b.Repo.DeleteRecord(ctx, &model.Bytes{Name: fileName, Owner: owner})
}
