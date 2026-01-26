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

type AuthHandler struct {
	rpc.UnimplementedAuthServiceServer
	Auth auth.Auther
}

func NewAuthHandler(auther auth.Auther) *AuthHandler {
	return &AuthHandler{Auth: auther}
}

func (k *AuthHandler) RegistUser(ctx context.Context, req *rpc.RegistUserRequest) (*rpc.RegistUserResponse, error) {
	token, err := k.Auth.Registration(ctx, req)

	// Ошибка создания пользователя.
	if err == errs.ErrCreateUser {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	// Ошибка уникального email.
	if err == errs.ErrEmailNotUnique {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err)
	}

	// Ошибки сервера.
	if err == nil || err == errs.ErrCreateToken {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Кладем токен в заголовок
	err = k.putDataToCtx(ctx, "authorization", token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	return &rpc.RegistUserResponse{Status: "ok"}, nil
}

func (k *AuthHandler) AuthUser(ctx context.Context, req *rpc.AuthUserRequest) (*rpc.AuthUserResponse, error) {
	token, err := k.Auth.Authentication(ctx, req)

	// Ошибка. Пользователь с таким email не найден.
	if err == errs.ErrUserNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	// Ошибка. Неверный пароль.
	if err == errs.ErrUserPassword {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err)
	}

	// Ошибки сервера.
	if err == errs.ErrCreateToken || err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	// Кладем токен в заголовок
	err = k.putDataToCtx(ctx, "authorization", token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	return &rpc.AuthUserResponse{Status: "ok"}, nil
}

func (k *AuthHandler) putDataToCtx(ctx context.Context, key, val string) error {
	return grpc.SetHeader(ctx, metadata.Pairs(key, val))
}
