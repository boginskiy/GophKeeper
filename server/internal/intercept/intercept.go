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

// return nil, status.Error(codes.Unauthenticated, "token is bad")
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

	// TODO! ...
	// Положить в хедер токен
	// Добавить данные в контекст
	// return handler(context.WithValue(ctx, auth.CtxUserID, UserID), req)

	return handler(ctx, req)
}

// Нужно предусмотреть, что чел идет как раз таки за регистрацией или аутентификацией

// Передаем токен и он валидный все ок. Разрешаем делать что дозволено
// Передаем токен и он невалидный или его нет. Значит пользователь проходит аутентификацию. Вводит пароль и почту.

// Если есть токен, то сначала передаем его на сервер.
// Если токен невалидный, то передаем данные для аутентификации (почту и пароль)
// Если пароль забыл, предусматриваем восстановление через почту или номер телефона
// Далее обновляем данные пароля.
