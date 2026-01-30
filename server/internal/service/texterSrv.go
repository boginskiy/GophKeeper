package service

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/repository"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
)

type TexterService struct {
	Cfg  config.Config
	Logg logg.Logger
	Repo repository.Repository[*model.Text]
}

func NewTexterService(config config.Config, logger logg.Logger, repo repository.Repository[*model.Text]) *TexterService {
	return &TexterService{
		Cfg:  config,
		Logg: logger,
		Repo: repo,
	}
}

func (t *TexterService) Create(ctx context.Context, req any) (any, error) {
	Req, ok := req.(*rpc.CreateRequest)
	if !ok {
		return nil, errs.ErrTypeConversion
	}
	return t.Repo.CreateRecord(model.NewText(Req))
}

func (t *TexterService) Read(ctx context.Context, req any) (any, error) {
	Req, ok := req.(*rpc.ReadRequest)
	if !ok {
		return nil, errs.ErrTypeConversion
	}
	return t.Repo.ReadRecord(&model.Text{Name: Req.Name, Owner: Req.Owner})
}

func (t *TexterService) ReadAll(ctx context.Context, req any) (any, error) {
	Req, ok := req.(*rpc.ReadAllRequest)
	if !ok {
		return nil, errs.ErrTypeConversion
	}
	return t.Repo.ReadAllRecord(&model.Text{Type: Req.Type, Owner: Req.Owner})
}

func (t *TexterService) Update(ctx context.Context, req any) (any, error) {
	Req, ok := req.(*rpc.CreateRequest)
	if !ok {
		return nil, errs.ErrTypeConversion
	}
	return t.Repo.UpdateRecord(&model.Text{Name: Req.Name, Tx: Req.Text, Owner: Req.Owner})
}

func (t *TexterService) Delete(ctx context.Context, req any) (any, error) {
	Req, ok := req.(*rpc.DeleteRequest)
	if !ok {
		return nil, errs.ErrTypeConversion
	}
	return t.Repo.DeleteRecord(&model.Text{Name: Req.Name, Owner: Req.Owner})
}
