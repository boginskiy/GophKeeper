package auth

import (
	"context"
)

// JWTer for a work of JWT authentication.
type JWTer interface {
	GetInfoAndValidJWT(token string) (*ExtraInfoToken, error)
	CreateJWT(*ExtraInfoToken) (string, error)
}

// Auther
type Auther[T any] interface {
	Identification(context.Context, any) (*ExtraInfoToken, bool)
	Registration(context.Context, T) (string, error)
	Authentication(context.Context, T) (string, error)
	Recovery(context.Context, T) (string, error)
}
