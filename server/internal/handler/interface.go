package handler

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/rpc"
)

type AuthService interface {
	Registration(context.Context, *rpc.RegistRequest) (*rpc.RegistResponse, error)
	Authentication(context.Context, *rpc.AuthRequest) (*rpc.AuthResponse, error)
	Recovery(context.Context, *rpc.RecovRequest) (*rpc.RecovResponse, error)
}
