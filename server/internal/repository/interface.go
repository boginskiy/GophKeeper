package repository

import "github.com/boginskiy/GophKeeper/server/internal/model"

type Repository interface {
}

type RepositoryUser interface {
	CreateRecord(*model.User) (*model.User, error)
	ReadRecord(email string) (*model.User, error)
}
