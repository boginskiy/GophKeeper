package auth

import (
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorHandler(err error) (bool, string) {
	if err == nil {
		return false, ""
	}

	// Ошибка локальная. Сервер не отвечает.
	if err == errs.ErrResponseServer {
		return true, "Server is unavailable, please try again later"
	}

	switch ModifyErrServerOnCode(err) {
	// Ошибка создания пользователя.
	case codes.InvalidArgument:
		return true, "User creation error"

	// Ошибка уникального email.
	case codes.AlreadyExists:
		return true, "Unique email error"

	// Ошибка создания токена или когда кладем токен в заголовок.
	case codes.Internal:
		return true, "An error in creating or transferring a token"

	// Ошибка запроса. Пользователь с таким email не найден.
	case codes.NotFound:
		return true, "User not found"

	// Ошибка запроса. Неверный пароль.
	case codes.Unauthenticated:
		return true, "Invalid password"

	// Неизвестная ошибка.
	default:
		return true, "Unknown error"
	}
}

func ModifyErrServerOnCode(err error) codes.Code {
	statusErr, ok := status.FromError(err)
	if !ok {
		return codes.OK
	}
	return statusErr.Code()
}
