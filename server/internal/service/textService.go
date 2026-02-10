package service

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/repo"
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

func (t *TextService) Create(ctx context.Context, mod *model.Text) (*model.Text, error) {
	return t.Repo.CreateRecord(ctx, mod)
}

func (t *TextService) Read(ctx context.Context, mod *model.Text) (*model.Text, error) {
	return t.Repo.ReadRecord(ctx, mod)
}

func (t *TextService) ReadAll(ctx context.Context, mod *model.Text) ([]*model.Text, error) {
	return t.Repo.ReadAllRecord(ctx, mod)
}

func (t *TextService) Update(ctx context.Context, mod *model.Text) (*model.Text, error) {
	return t.Repo.UpdateRecord(ctx, mod)
}

func (t *TextService) Delete(ctx context.Context, mod *model.Text) (*model.Text, error) {
	return t.Repo.DeleteRecord(ctx, mod)
}
