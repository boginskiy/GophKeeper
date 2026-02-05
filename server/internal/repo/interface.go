package repo

import "context"

type CreateReader[T any] interface {
	CreateRecord(ctx context.Context, obj T) (T, error)
	ReadRecord(ctx context.Context, obj T) (T, error)
}

type Repository[T any] interface {
	CreateReader[T]
	ReadAllRecord(ctx context.Context, obj T) ([]T, error)
	UpdateRecord(ctx context.Context, obj T) (T, error)
	DeleteRecord(ctx context.Context, obj T) (T, error)
}
