package repo

type CreateReader[T any] interface {
	CreateRecord(obj T) (T, error)
	ReadRecord(obj T) (T, error)
}

type Repository[T any] interface {
	CreateReader[T]
	ReadAllRecord(obj T) ([]T, error)
	UpdateRecord(obj T) (T, error)
	DeleteRecord(obj T) (T, error)
}
