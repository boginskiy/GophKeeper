package service

import (
	"bufio"
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

type ByterService struct {
	Cfg         config.Config
	Logg        logg.Logger
	Repo        repo.CreateReader[*model.Bytes]
	FileHdler   utils.FileHandler
	FileManager manager.FileManager
}

func NewByterService(
	config config.Config,
	logger logg.Logger,
	repo repo.CreateReader[*model.Bytes],
	fileHdler utils.FileHandler,
	fileManager manager.FileManager,
) *ByterService {

	return &ByterService{
		Cfg:         config,
		Logg:        logger,
		Repo:        repo,
		FileHdler:   fileHdler,
		FileManager: fileManager,
	}
}

func (b *ByterService) Upload(stream any) (any, error) {
	Stm, ok := stream.(rpc.ByterService_UploadServer)
	if !ok {
		return nil, errs.ErrTypeConversion
	}

	modBytes := &model.Bytes{}

	// insert FileSize, FileName, FileOwner in modBytes
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

func (b *ByterService) uploadStream(stream rpc.ByterService_UploadServer, modBytes *model.Bytes) (int64, error) {
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

func (b *ByterService) Unload(stream any) (any, error) {
	Stm, okS := stream.(rpc.ByterService_UnloadServer)
	if !okS {
		return nil, errs.ErrTypeConversion
	}

	fileName, errFile := manager.TakeClientValueFromCtx(Stm.Context(), "file_name", 0)
	owner, errOwner := manager.TakeServerValueFromCtx(Stm.Context(), manager.EmailCtx)

	if errFile != nil || errOwner != nil {
		return nil, utils.DefinErr(errFile, errOwner)
	}

	// Take record from DB.
	modBytes, err := b.Repo.ReadRecord(&model.Bytes{Name: fileName, Owner: owner})
	if err != nil {
		return nil, err
	}

	// Check and Read file
	if !b.FileHdler.CheckOfFile(modBytes.Path) {
		return nil, errs.ErrFileNotFound
	}

	file, err := b.FileHdler.ReadOrCreateFile(modBytes.Path, manager.MOD)
	if err != nil {
		return nil, err
	}

	defer modBytes.Descr.Close()
	modBytes.Descr = file

	err = b.unloadStream(Stm, modBytes)
	if err != nil {
		return nil, err
	}

	return modBytes, nil
}

func (b *ByterService) unloadStream(stream rpc.ByterService_UnloadServer, modBytes *model.Bytes) error {
	// Buffer 1KB.
	buffer := make([]byte, 1024)

	// Run stream.
	for {
		n, err := modBytes.Descr.Read(buffer)

		if err != nil && err != io.EOF {
			return errs.ErrReadFileToBuff.Wrap(err)
		}
		if n == 0 {
			break
		}
		// Part of request.
		chunk := &rpc.UnloadBytesResponse{Content: buffer[:n]}

		// Send part of request.
		if err := stream.Send(chunk); err != nil {
			return errs.ErrSendChankFile.Wrap(err)
		}
	}
	return nil
}
