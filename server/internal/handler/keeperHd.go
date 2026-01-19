package handler

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/auth"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type KeeperHandler struct {
	rpc.UnimplementedKeeperServiceServer
	// Add Service
	Auth auth.Auther
}

func NewKeeperHandler(auther auth.Auther) *KeeperHandler {
	return &KeeperHandler{Auth: auther}
}

func (k *KeeperHandler) RegistUser(ctx context.Context, req *rpc.RegistUserRequest) (*rpc.RegistUserResponse, error) {
	token, err := k.Auth.Registration(ctx, req)

	if err == errs.ErrCreateUser {
		// Ошибка создания пользователя.
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	if err == errs.ErrEmailNotUnique {
		// Ошибка уникального email.
		return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
	}

	if err == errs.ErrCreateToken {
		// Ошибка создания токена.
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	// Кладем токен в заголовок
	err = k.putDataToCtx(ctx, "authorization", token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &rpc.RegistUserResponse{Status: "ok"}, nil
}

func (k *KeeperHandler) AuthUser(ctx context.Context, req *rpc.AuthUserRequest) (*rpc.AuthUserResponse, error) {
	token, err := k.Auth.Authentication(ctx, req)

	if err == errs.ErrUserNotFound {
		// Ошибка. Пользователь с таким email не найден.
		return nil, status.Errorf(codes.NotFound, "%s", err.Error())
	}

	if err == errs.ErrUserPassword {
		// Ошибка. Неверный пароль.
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}

	if err == errs.ErrCreateToken {
		// Ошибка создания токена.
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	// Кладем токен в заголовок
	err = k.putDataToCtx(ctx, "authorization", token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &rpc.AuthUserResponse{Status: "ok"}, nil
}

func (k *KeeperHandler) putDataToCtx(ctx context.Context, key, val string) error {
	return grpc.SetHeader(ctx, metadata.Pairs(key, val))
}
