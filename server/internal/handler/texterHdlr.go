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
	Service service.Servicer
}

func NewTexterHandler(srv service.Servicer) *TexterHandler {
	return &TexterHandler{Service: srv}
}

func (k *TexterHandler) Create(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
	obj, err := k.Service.Create(ctx, req)

	// Тут надо сделать проверку на уникальность типа записи.
	if err == errs.ErrEmailNotUnique {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err)
	}

	// Иные ошибки сервера.
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Conversion.
	text, ok := obj.(*model.Text)
	if !ok {
		return nil, status.Errorf(codes.Internal, "%s", errs.ErrTypeConversion)
	}

	return &rpc.CreateResponse{
		Status:    "created",
		UpdatedAt: utils.ConvertDtStr(text.UpdatedAt)}, nil
}

func (k *TexterHandler) Read(ctx context.Context, req *rpc.ReadRequest) (*rpc.ReadResponse, error) {
	obj, err := k.Service.Read(ctx, req)

	// Данные не найдены.
	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	// Иные ошибки сервера.
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Conversion.
	text, ok := obj.(*model.Text)
	if !ok {
		return nil, status.Errorf(codes.Internal, "%s", errs.ErrTypeConversion)
	}

	return &rpc.ReadResponse{
		Status:    "read",
		Name:      text.Name,
		Text:      text.Tx,
		UpdatedAt: utils.ConvertDtStr(text.UpdatedAt)}, nil
}

func (k *TexterHandler) ReadAll(ctx context.Context, req *rpc.ReadAllRequest) (*rpc.ReadAllResponse, error) {
	objs, err := k.Service.ReadAll(ctx, req)

	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	texts, ok := objs.([]*model.Text)
	if !ok {
		return nil, status.Errorf(codes.Internal, "%s", errs.ErrTypeConversion)
	}

	textResponses := make([]*rpc.TextResponse, len(texts))
	for i, text := range texts {
		textResponses[i] = &rpc.TextResponse{
			Name:      text.Name,
			Text:      text.Tx,
			UpdatedAt: utils.ConvertDtStr(text.UpdatedAt),
		}
	}

	return &rpc.ReadAllResponse{
		Status:        "read",
		TextResponses: textResponses}, nil
}

func (k *TexterHandler) Update(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
	obj, err := k.Service.Update(ctx, req)

	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	text, ok := obj.(*model.Text)
	if !ok {
		return nil, status.Errorf(codes.Internal, "%s", errs.ErrTypeConversion)
	}

	return &rpc.CreateResponse{
		Status:    "update",
		UpdatedAt: utils.ConvertDtStr(text.UpdatedAt)}, nil
}

func (k *TexterHandler) Delete(ctx context.Context, req *rpc.DeleteRequest) (*rpc.DeleteResponse, error) {
	obj, err := k.Service.Delete(ctx, req)

	if err == errs.ErrDataNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	text, ok := obj.(*model.Text)
	if !ok {
		return nil, status.Errorf(codes.Internal, "%s", errs.ErrTypeConversion)
	}

	return &rpc.DeleteResponse{Status: "delete", Name: text.Name}, nil
}
