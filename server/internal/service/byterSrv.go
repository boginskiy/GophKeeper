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
	Repo        repo.Repository[*model.Bytes]
	FileHdler   utils.FileHandler
	FileManager manager.FileManager
}

func NewByterService(
	config config.Config,
	logger logg.Logger,
	repo repo.Repository[*model.Bytes],
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
	return nil, nil
}

// type Bytes struct {
// 	Name      string
// 	Path      string
// 	Descr     *os.File
// 	Size      int64
// 	Type      string
// 	Owner     string
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }
