package service

import (
	"io"
	"time"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/infra"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/repo"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
)

type UnloadService struct {
	Cfg         config.Config
	Logger      logg.Logger
	FileHandler utils.FileHandler
	Repo        repo.RepoCreateReader[*model.Bytes]
}

func NewUnloadService(
	config config.Config,
	logger logg.Logger,
	fileHandler utils.FileHandler,
	repo repo.RepoCreateReader[*model.Bytes]) *UnloadService {

	return &UnloadService{
		Cfg:         config,
		Logger:      logger,
		FileHandler: fileHandler,
		Repo:        repo,
	}
}

func (s *UnloadService) Load(stream rpc.ByterService_UnloadServer, modBytes *model.Bytes) (*model.Bytes, error) {
	// Кладем данные в заголовок для клиента
	errUp := infra.PutDataToCtx(stream.Context(), "updated_at", utils.ConvertDtStr(utils.ConversDtToTableView(time.Now())))
	errSz := infra.PutDataToCtx(stream.Context(), "sent_size", modBytes.ReceivedSize)
	errFn := infra.PutDataToCtx(stream.Context(), "file_name", modBytes.Name)

	if errUp != nil || errSz != nil || errFn != nil {
		return nil, utils.DefinErr(errUp, errSz, errFn)
	}

	// Check and Read file
	if !s.FileHandler.CheckOfFile(modBytes.Path) {
		return nil, errs.ErrFileNotFound
	}

	file, err := s.FileHandler.ReadOrCreateFile(modBytes.Path, infra.MOD)
	if err != nil {
		return nil, err
	}

	defer modBytes.Descr.Close()
	modBytes.Descr = file

	err = s.unloadStream(stream, modBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *UnloadService) unloadStream(stream rpc.ByterService_UnloadServer, modBytes *model.Bytes) error {
	// Buffer 1KB.
	buffer := make([]byte, 1<<10)

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
