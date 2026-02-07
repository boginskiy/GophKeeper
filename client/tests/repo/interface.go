package repo

type RepoCreateReader[T any] interface {
	CreateRecord(obj T) (T, error)
	ReadRecord(obj T) (T, error)
}

type Repository[T any] interface {
	RepoCreateReader[T]
	ReadAllRecord(obj T) ([]T, error)
	UpdateRecord(obj T) (T, error)
	DeleteRecord(obj T) (T, error)
}
