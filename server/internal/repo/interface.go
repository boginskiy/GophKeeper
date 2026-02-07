package repo

import "context"

type Repository[T any] interface {
	RepoCreateReader[T]
	ReadAllRecord(ctx context.Context, obj T) ([]T, error)
	UpdateRecord(ctx context.Context, obj T) (T, error)
	DeleteRecord(ctx context.Context, obj T) (T, error)
}

type RepoCreateReadUpdater[T any] interface {
	RepoCreateReader[T]
	RepoUpdater[T]
}

type RepoUpdater[T any] interface {
	UpdateRecord(ctx context.Context, obj T) (T, error)
}

type RepoCreateReader[T any] interface {
	CreateRecord(ctx context.Context, obj T) (T, error)
	ReadRecord(ctx context.Context, obj T) (T, error)
}
