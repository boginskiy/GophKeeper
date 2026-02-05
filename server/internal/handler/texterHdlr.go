package handler

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/internal/service"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TexterHandler struct {
	rpc.UnimplementedTexterServiceServer
	Service service.TextServicer[*model.Text]
}

func NewTexterHandler(srv service.TextServicer[*model.Text]) *TexterHandler {
	return &TexterHandler{Service: srv}
}

func (k *TexterHandler) Create(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
	modText, err := k.Service.Create(ctx, req)

	// Тут надо сделать проверку на уникальность типа записи.
	if err == errs.ErrEmailNotUnique {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err)
	}

	// Иные ошибки сервера.
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	return &rpc.CreateResponse{
		Status:    "created",
		UpdatedAt: utils.ConvertDtStr(modText.UpdatedAt)}, nil
}

func (k *TexterHandler) Read(ctx context.Context, req *rpc.ReadRequest) (*rpc.ReadResponse, error) {
	modText, err := k.Service.Read(ctx, req)

	// Данные не найдены.
	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	// Иные ошибки сервера.
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	return &rpc.ReadResponse{
		Status:    "read",
		Name:      modText.Name,
		Text:      modText.Content,
		UpdatedAt: utils.ConvertDtStr(modText.UpdatedAt)}, nil
}

func (k *TexterHandler) ReadAll(ctx context.Context, req *rpc.ReadAllRequest) (*rpc.ReadAllResponse, error) {
	modTexts, err := k.Service.ReadAll(ctx, req)

	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	textResponses := make([]*rpc.TextResponse, len(modTexts))
	for i, text := range modTexts {
		textResponses[i] = &rpc.TextResponse{
			Name:      text.Name,
			Text:      text.Content,
			UpdatedAt: utils.ConvertDtStr(text.UpdatedAt),
		}
	}

	return &rpc.ReadAllResponse{
		Status:        "read",
		TextResponses: textResponses}, nil
}

func (k *TexterHandler) Update(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
	modText, err := k.Service.Update(ctx, req)

	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	return &rpc.CreateResponse{
		Status:    "update",
		UpdatedAt: utils.ConvertDtStr(modText.UpdatedAt)}, nil
}

func (k *TexterHandler) Delete(ctx context.Context, req *rpc.DeleteRequest) (*rpc.DeleteResponse, error) {
	modText, err := k.Service.Delete(ctx, req)

	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	return &rpc.DeleteResponse{Status: "delete", Name: modText.Name}, nil
}
