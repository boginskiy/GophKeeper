package service

import (
	"fmt"
	"io"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/auth"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/repository"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
)

type ByterService struct {
	Cfg  config.Config
	Logg logg.Logger
	Repo repository.Repository[*model.Text]
}

func NewByterService(config config.Config, logger logg.Logger, repo repository.Repository[*model.Text]) *ByterService {
	return &ByterService{
		Cfg:  config,
		Logg: logger,
		Repo: repo,
	}
}

func (b *ByterService) Upload(stream any) (any, error) {
	Stm, ok := stream.(rpc.ByterService_UploadServer)
	if !ok {
		return nil, errs.ErrTypeConversion
	}

	owner, ok := Stm.Context().Value(auth.EmailCtx).(string)
	if !ok {
		return nil, errs.ErrDataOwner
	}

	bytes := make([]byte, 0, 1024)

	// Обработка запроса.
	for {
		req, err := Stm.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		bytes = req.Content

		// Как лучше писать сохранять данные ? bytes

	}

	//
	fmt.Println(owner, bytes)

	return &rpc.UploadBytesResponse{}, nil
}

func (b *ByterService) Unload(stream any) (any, error) {

	return nil, nil
}

// type ServicerByter interface {
// 	Upload(context.Context, any) (any, error)
// 	Unload(context.Context, any) (any, error)
// }
