package repository

import (
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/model"
)

type RepoUser struct {
	Store map[string]*model.User
}

func NewRepoUser() *RepoUser {
	return &RepoUser{Store: make(map[string]*model.User, 10)}
}

func (u *RepoUser) CreateRecord(user *model.User) (*model.User, error) {
	if _, ok := u.Store[user.Email]; ok {
		return nil, errs.ErrEmailNotUnique
	}
	user.ID = len(u.Store) + 1
	u.Store[user.Email] = user
	return user, nil
}

func (u *RepoUser) ReadRecord(email string) (*model.User, error) {
	user, ok := u.Store[email]
	if !ok {
		return nil, errs.ErrUserNotFound
	}
	return user, nil
}
