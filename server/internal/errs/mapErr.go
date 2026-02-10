package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrMapper interface {
	Mapping(err error) error
}

type ErrMap struct {
}

func NewErrMap() *ErrMap {
	return &ErrMap{}
}

func (e *ErrMap) Mapping(err error) error {
	if err == nil {
		return nil
	}

	switch err {

	// Ошибка создания пользователя.
	case ErrCreateUser:
		return status.Errorf(codes.InvalidArgument, "%s", err)

	// Ошибка уникального email.
	case ErrEmailNotUnique:
		return status.Errorf(codes.AlreadyExists, "%s", err)

	// Ошибка. Пользователь с таким email не найден.
	case ErrUserNotFound, ErrFileNotFound:
		return status.Errorf(codes.NotFound, "%s", err)

	// Ошибка. Неверный пароль.
	case ErrUserPassword:
		return status.Errorf(codes.Unauthenticated, "%s", err)

	// Ошибки сервера.
	default:
		return status.Errorf(codes.Internal, "%s", err)
	}
}
