package service

import (
	"bufio"
	"io"
	"strconv"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/infra"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/repo"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
	"github.com/boginskiy/GophKeeper/server/pkg"
)

type UploadService struct {
	Cfg           config.Config
	Logger        logg.Logger
	FileHandler   utils.FileHandler
	FileService   infra.Filer
	CryptoService pkg.Crypter
	Repo          repo.RepoCreateReader[*model.Bytes]
}

func NewUploadService(
	config config.Config,
	logger logg.Logger,
	fileHandler utils.FileHandler,
	fileService infra.Filer,
	cryptoService pkg.Crypter,
	repo repo.RepoCreateReader[*model.Bytes]) *UploadService {

	tmp := &UploadService{
		Cfg:           config,
		Logger:        logger,
		FileHandler:   fileHandler,
		FileService:   fileService,
		CryptoService: cryptoService,
		Repo:          repo,
	}

	// Start CryptoService.
	tmp.CryptoService.Start([]byte("CryptoKey"))
	return tmp
}

func (s *UploadService) Load(stream rpc.ByterService_UploadServer, modBytes *model.Bytes) (*model.Bytes, error) {
	defer modBytes.Descr.Close()

	cnt, err := s.uploadStream(stream, modBytes)
	if err != nil {
		return nil, errs.ErrRunStream.Wrap(err)
	}

	modBytes.ReceivedSize = strconv.FormatInt(cnt, 10)
	return s.Repo.CreateRecord(stream.Context(), modBytes)
}

func (s *UploadService) uploadStream(stream rpc.ByterService_UploadServer, modBytes *model.Bytes) (int64, error) {
	// Writer
	writer := bufio.NewWriter(modBytes.Descr)
	s.CryptoService.Reset()

	var ClientSignature []byte
	var CNT int64

	for {
		// Обработка запроса.
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		// Check of Crypto signature.
		if len(req.CryptoSignature) > 0 {
			ClientSignature = req.CryptoSignature
		}

		// Content loading to file.
		nn, err := writer.Write(req.Content)
		s.CryptoService.Write(req.Content)

		if err != nil {
			return 0, err
		}

		CNT += int64(nn)
	}

	// Проверка цифровой подписи.
	ok := s.CryptoService.CheckSignature(ClientSignature)
	if !ok {
		return 0, pkg.ErrCheckCryptoSignature
	}

	err := writer.Flush()
	if err != nil {
		return 0, err
	}

	return CNT, nil
}
