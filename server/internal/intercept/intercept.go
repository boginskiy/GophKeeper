package intercept

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/auth"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Interceptor interface {
	WithAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

type Intercept struct {
	Cfg  config.Config
	Logg logg.Logger
	Auth auth.Auther
}

func NewIntercept(config config.Config, logger logg.Logger, auther auth.Auther) *Intercept {
	return &Intercept{
		Cfg:  config,
		Logg: logger,
		Auth: auther,
	}
}

func (i *Intercept) WithAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	infoToken, ok := i.Auth.Identification(ctx, req)

	// Bad identification. Send response about Authentication.
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "token is bad")
	}

	// Request to Authentication or Registration.
	if ok && infoToken == nil {
		handler(ctx, req)
	}

	// Good identification.

	ctx = context.WithValue(ctx, auth.PhoneCtx, infoToken.PhoneNumber)
	ctx = context.WithValue(ctx, auth.EmailCtx, infoToken.Email)

	return handler(ctx, req)
}
