package service

import (
	"bufio"
	"context"
	"io"
	"strconv"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/manager"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/repo"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
)

type BytesService struct {
	Cfg         config.Config
	Logg        logg.Logger
	Repo        repo.Repository[*model.Bytes]
	FileHdler   utils.FileHandler
	FileManager manager.FileManager
}

func NewBytesService(
	config config.Config,
	logger logg.Logger,
	repo repo.Repository[*model.Bytes],
	fileHdler utils.FileHandler,
	fileManager manager.FileManager,
) *BytesService {

	return &BytesService{
		Cfg:         config,
		Logg:        logger,
		Repo:        repo,
		FileHdler:   fileHdler,
		FileManager: fileManager,
	}
}

func (b *BytesService) Upload(stream any) (*model.Bytes, error) {
	Stm, ok := stream.(rpc.ByterService_UploadServer)
	if !ok {
		return nil, errs.ErrTypeConversion
	}

	modBytes := &model.Bytes{}

	// insert FileSize, DataType, FileName, FileOwner in modBytes
	err := modBytes.InsertValuesFromCtx(Stm.Context())
	if err != nil {
		return nil, errs.ErrDataCtx // Ошибка запроса request клиента
	}

	// File for data saving
	file, path, err := b.FileManager.CreateFileInStore(modBytes)
	if err != nil {
		return nil, errs.ErrCreateFile.Wrap(err)
	}

	modBytes.Descr, modBytes.Path = file, path
	defer file.Close()

	cnt, err := b.uploadStream(Stm, modBytes)
	if err != nil {
		return nil, errs.ErrRunStream.Wrap(err)
	}

	modBytes.ReceivedSize = strconv.FormatInt(cnt, 10)
	return b.Repo.CreateRecord(modBytes)
}

func (b *BytesService) uploadStream(stream rpc.ByterService_UploadServer, modBytes *model.Bytes) (int64, error) {
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

func (b *BytesService) Read(ctx context.Context, req any) (*model.Bytes, error) {
	// Info from context.
	fileName, err := manager.TakeClientValueFromCtx(ctx, "file_name", 0)
	if err != nil {
		return nil, err
	}

	owner, err := manager.TakeServerValueFromCtx(ctx, manager.EmailCtx)
	if err != nil {
		return nil, err
	}

	return b.Repo.ReadRecord(&model.Bytes{Name: fileName, Owner: owner})
}

func (b *BytesService) ReadAll(ctx context.Context, req any) ([]*model.Bytes, error) {
	// Info from context.
	fileType, err := manager.TakeClientValueFromCtx(ctx, "file_type", 0)
	if err != nil {
		return nil, err
	}

	owner, err := manager.TakeServerValueFromCtx(ctx, manager.EmailCtx)
	if err != nil {
		return nil, err
	}

	return b.Repo.ReadAllRecord(&model.Bytes{Type: fileType, Owner: owner})
}

func (b *BytesService) Delete(ctx context.Context, req any) (*model.Bytes, error) {
	return nil, nil
}
