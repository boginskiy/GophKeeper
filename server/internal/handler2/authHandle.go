package handler2

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/auth"
	"github.com/boginskiy/GophKeeper/server/internal/infra"
	"github.com/boginskiy/GophKeeper/server/internal/model"
)

type AuthHandle struct {
	AuthService auth.Auther[*model.User]
}

func NewAuthHandle(service auth.Auther[*model.User]) *AuthHandle {
	return &AuthHandle{AuthService: service}
}

func (k *AuthHandle) HandleRegistration(ctx context.Context, req *model.User) (string, error) {
	token, err := k.AuthService.Registration(ctx, req)
	if err != nil {
		return "", err
	}

	// Put token to header.
	err = infra.PutDataToCtx(ctx, "authorization", token)
	if err != nil {
		return "", err
	}

	return "ok", nil
}

func (k *AuthHandle) HandleAuthentication(ctx context.Context, req *model.User) (string, error) {
	token, err := k.AuthService.Authentication(ctx, req)
	if err != nil {
		return "", err
	}

	// Put token to header.
	err = infra.PutDataToCtx(ctx, "authorization", token)
	if err != nil {
		return "", err
	}

	return "ok", nil
}

func (k *AuthHandle) HandleRecovery(ctx context.Context, req *model.User) (string, error) {
	token, err := k.AuthService.Recovery(ctx, req)
	if err != nil {
		return "", err
	}

	// Put token to header.
	err = infra.PutDataToCtx(ctx, "authorization", token)
	if err != nil {
		return "", err
	}

	return "ok", nil
}
