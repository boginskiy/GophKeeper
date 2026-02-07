package repo

import (
	"context"
	"database/sql"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/db"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
)

type RepoUser struct {
	Cfg    config.Config
	Logg   logg.Logger
	Store  db.DataBase[*sql.DB]
	SqlDB  *sql.DB
	NameTb string
}

func NewRepoUser(cfg config.Config, logger logg.Logger, db db.DataBase[*sql.DB]) *RepoUser {
	return &RepoUser{
		Cfg:    cfg,
		Logg:   logger,
		Store:  db,
		SqlDB:  db.GetDB(),
		NameTb: "users",
	}
}

func (u *RepoUser) UpdateRecord(ctx context.Context, user *model.User) (*model.User, error) {
	row := u.SqlDB.QueryRowContext(ctx,
		`UPDATE users
		 SET password = $1
		 WHERE email = $2
		 RETURNING id, phone_number;`,
		user.Password, user.Email)

	return user, row.Scan(&user.ID, &user.PhoneNumber)
}

func (u *RepoUser) CreateRecord(ctx context.Context, user *model.User) (*model.User, error) {
	row := u.SqlDB.QueryRowContext(ctx,
		`INSERT INTO users (name, email, phone_number, password)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id;`,
		user.UserName, user.Email, user.PhoneNumber, user.Password)

	return user, row.Scan(&user.ID)
}

func (u *RepoUser) ReadRecord(ctx context.Context, user *model.User) (*model.User, error) {
	row := u.SqlDB.QueryRowContext(ctx,
		`SELECT id, name, password, phone_number, role
		 FROM users 
		 WHERE email = $1;`,
		user.Email)

	if err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Password,
		&user.PhoneNumber,
		&user.Role); err != nil {
		return nil, err
	}
	return user, nil
}
