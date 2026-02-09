package service

import "context"

type TextServicer[T any] interface {
	Create(context.Context, any) (T, error)
	Read(context.Context, any) (T, error)
	ReadAll(context.Context, any) ([]T, error)
	Update(context.Context, any) (T, error)
	Delete(context.Context, any) (T, error)
}

type BytesServicer[T any] interface {
	Read(context.Context, any) (T, error)
	ReadAll(context.Context, any) ([]T, error)
	Delete(context.Context, any) (T, error)

	// Убрать в другой сервис
	Upload(any) (T, error)
}

type LoadServicer[ST, M any] interface {
	Prepar(stream ST) (model M, err error)
	Load(stream ST, model M) (M, error)
}
