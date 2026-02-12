package handler

import (
	"context"
)

type AuthHandler[T any] interface {
	HandleAuthentication(context.Context, T) (string, error)
	HandleRegistration(context.Context, T) (string, error)
	HandleRecovery(context.Context, T) (string, error)
}

type ByteHandler[T any] interface {
	HandleDelete(context.Context, T) (string, error)
	HandleReadAll(context.Context, T) ([]T, error)
	HandleRead(context.Context, T) (T, error)

	// HandleUpload(context.Context, T) (string, error)
	// HandleUnload(context.Context, T) (string, error)
}

type TextHandler[T any] interface {
	HandleCreate(context.Context, T) (T, error)
	HandleDelete(context.Context, T) (T, error)
	HandleRead(context.Context, T) (T, error)
	HandleReadAll(context.Context, T) ([]T, error)
	HandleUpdate(context.Context, T) (T, error)
}
