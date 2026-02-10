package handler2

import "context"

type AuthHandler[T any] interface {
	HandleAuthentication(context.Context, T) (string, error)
	HandleRegistration(context.Context, T) (string, error)
	HandleRecovery(context.Context, T) (string, error)
}

type ByteHandler[T any] interface {
	HandleReadAll(context.Context, T) ([]T, error)
	// HandleUpload(context.Context, T) (string, error)
	// HandleUnload(context.Context, T) (string, error)
	HandleDelete(context.Context, T) (string, error)
	HandleRead(context.Context, T) (T, error)
}
