package auth

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/rpc"
)

// JWTer for a work of JWT authentication.
type JWTer interface {
	GetInfoAndValidJWT(token string) (*ExtraInfoToken, error)
	CreateJWT(*ExtraInfoToken) (string, error)
}

// Auther .
type Auther interface {
	Identification(context.Context, any) (*ExtraInfoToken, bool)
	Registration(context.Context, *rpc.RegistUserRequest) (string, error)
	Authentication(context.Context, *rpc.AuthUserRequest) (string, error)
}
