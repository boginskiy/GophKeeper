package service

import (
	"io"

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
	Cfg       config.Config
	Logg      logg.Logger
	FileHdler utils.FileHandler
	Repo      repo.CreateReader[*model.Bytes]
}

func NewUnloadService(
	config config.Config,
	logger logg.Logger,
	fileHdler utils.FileHandler,
	repo repo.CreateReader[*model.Bytes]) *UnloadService {

	return &UnloadService{
		Cfg:       config,
		Logg:      logger,
		FileHdler: fileHdler,
		Repo:      repo,
	}
}

func (s *UnloadService) Prepar(stream rpc.ByterService_UnloadServer) (*model.Bytes, error) {
	// Info from context.
	fileName, err := infra.TakeClientValueFromCtx(stream.Context(), "file_name", 0)
	if err != nil {
		return nil, err
	}

	owner, err := infra.TakeServerValueFromCtx(stream.Context(), infra.EmailCtx)
	if err != nil {
		return nil, err
	}

	// Take record from DataBase.
	modBytes, err := s.Repo.ReadRecord(&model.Bytes{Name: fileName, Owner: owner})
	if err != nil {
		return nil, err
	}

	return modBytes, nil
}

func (s *UnloadService) Load(stream rpc.ByterService_UnloadServer, modBytes *model.Bytes) error {
	// Check and Read file
	if !s.FileHdler.CheckOfFile(modBytes.Path) {
		return errs.ErrFileNotFound
	}

	file, err := s.FileHdler.ReadOrCreateFile(modBytes.Path, infra.MOD)
	if err != nil {
		return err
	}

	defer modBytes.Descr.Close()
	modBytes.Descr = file

	err = s.unloadStream(stream, modBytes)
	if err != nil {
		return err
	}

	return nil
}

func (s *UnloadService) unloadStream(stream rpc.ByterService_UnloadServer, modBytes *model.Bytes) error {
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
