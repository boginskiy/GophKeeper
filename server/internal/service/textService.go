package service

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/repo"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
)

type TextService struct {
	Cfg    config.Config
	Logger logg.Logger
	Repo   repo.Repository[*model.Text]
}

func NewTextService(config config.Config, logger logg.Logger, repo repo.Repository[*model.Text]) *TextService {
	return &TextService{
		Cfg:    config,
		Logger: logger,
		Repo:   repo,
	}
}

func (t *TextService) Create(ctx context.Context, req any) (*model.Text, error) {
	Req, ok := req.(*rpc.CreateRequest)
	if !ok {
		return nil, errs.ErrTypeConversion
	}
	// userID is in ctx
	return t.Repo.CreateRecord(ctx, model.NewText(Req))
}

func (t *TextService) Read(ctx context.Context, req any) (*model.Text, error) {
	Req, ok := req.(*rpc.ReadRequest)
	if !ok {
		return nil, errs.ErrTypeConversion
	}
	return t.Repo.ReadRecord(ctx, &model.Text{Name: Req.Name, Owner: Req.Owner})
}

func (t *TextService) ReadAll(ctx context.Context, req any) ([]*model.Text, error) {
	Req, ok := req.(*rpc.ReadAllRequest)
	if !ok {
		return nil, errs.ErrTypeConversion
	}
	return t.Repo.ReadAllRecord(ctx, &model.Text{Type: Req.Type, Owner: Req.Owner})
}

func (t *TextService) Update(ctx context.Context, req any) (*model.Text, error) {
	Req, ok := req.(*rpc.CreateRequest)
	if !ok {
		return nil, errs.ErrTypeConversion
	}
	return t.Repo.UpdateRecord(ctx, &model.Text{Name: Req.Name, Content: Req.Text, Owner: Req.Owner})
}

func (t *TextService) Delete(ctx context.Context, req any) (*model.Text, error) {
	Req, ok := req.(*rpc.DeleteRequest)
	if !ok {
		return nil, errs.ErrTypeConversion
	}
	return t.Repo.DeleteRecord(ctx, &model.Text{Name: Req.Name, Owner: Req.Owner})
}
