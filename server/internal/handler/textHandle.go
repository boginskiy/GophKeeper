package handler

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/service"
)

type TextHandle struct {
	TextService service.TextServicer[*model.Text]
}

func NewTextHandle(
	textService service.TextServicer[*model.Text]) *TextHandle {
	return &TextHandle{TextService: textService}
}

func (t *TextHandle) HandleCreate(ctx context.Context, req *model.Text) (*model.Text, error) {
	return t.TextService.Create(ctx, req)
}

func (t *TextHandle) HandleRead(ctx context.Context, req *model.Text) (*model.Text, error) {
	return t.TextService.Read(ctx, req)
}

func (t *TextHandle) HandleReadAll(ctx context.Context, req *model.Text) ([]*model.Text, error) {
	return t.TextService.ReadAll(ctx, req)
}

func (t *TextHandle) HandleUpdate(ctx context.Context, req *model.Text) (*model.Text, error) {
	return t.TextService.Update(ctx, req)
}

func (t *TextHandle) HandleDelete(ctx context.Context, req *model.Text) (*model.Text, error) {
	return t.TextService.Delete(ctx, req)
}
