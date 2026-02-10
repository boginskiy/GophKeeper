package service

import "context"

type TextServicer[T any] interface {
	Create(context.Context, T) (T, error)
	Read(context.Context, T) (T, error)
	ReadAll(context.Context, T) ([]T, error)
	Update(context.Context, T) (T, error)
	Delete(context.Context, T) (T, error)
}

type BytesServicer[T any] interface {
	Read(context.Context, T) (T, error)
	ReadAll(context.Context, T) ([]T, error)
	Delete(context.Context, T) (T, error)
}

type LoadServicer[ST, M any] interface {
	Load(stream ST, model M) (M, error)
}
