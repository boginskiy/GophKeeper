package intercept

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/auth"
	"github.com/boginskiy/GophKeeper/server/internal/infra"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServIntercept struct {
	Cfg    config.Config
	Logger logg.Logger
	Auth   auth.Auther[*model.User]
}

func NewServIntercept(config config.Config, logger logg.Logger, auther auth.Auther[*model.User]) *ServIntercept {
	return &ServIntercept{
		Cfg:    config,
		Logger: logger,
		Auth:   auther,
	}
}

func (i *ServIntercept) ServAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	infoToken, ok := i.Auth.Identification(ctx, req)

	// Bad identification. Send response about Authentication.
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "token is bad")
	}

	// Request to Authentication or Registration.
	if ok && infoToken == nil {
		return handler(ctx, req)
	}
	// Good identification.
	ctx = context.WithValue(ctx, infra.PhoneCtx, infoToken.PhoneNumber)
	ctx = context.WithValue(ctx, infra.EmailCtx, infoToken.Email)
	ctx = context.WithValue(ctx, infra.IDCtx, infoToken.ID)

	return handler(ctx, req)
}

func (i *ServIntercept) StreamServAuth(service interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// На этом этапе мы через одиночный запрос в ServAuth зарегистрировались/аутентифицировались.
	// В данном перехватчике мы делаем сверку токен и кладем в контекст доп инфу.

	ctx := ss.Context()

	infoToken, ok := i.Auth.Identification(ctx, ss)
	if !ok {
		return status.Error(codes.Unauthenticated, "token is bad")
	}

	// New context.
	ctx = context.WithValue(ctx, infra.PhoneCtx, infoToken.PhoneNumber)
	ctx = context.WithValue(ctx, infra.EmailCtx, infoToken.Email)
	ctx = context.WithValue(ctx, infra.IDCtx, infoToken.ID)

	// Оболочка ServerStream с новым контекстом
	wrapSS := &WrapServerStream{
		ServerStream: ss,
		Ctx:          ctx,
	}

	return handler(service, wrapSS)
}
