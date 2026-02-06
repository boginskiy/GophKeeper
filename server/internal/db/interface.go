package db

type DataBase[T any] interface {
	CheckOpen() bool
	GetDB() T
	CloseDB()
}
