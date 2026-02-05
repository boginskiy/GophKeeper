package repo

import (
	"context"
	"database/sql"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/db"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
)

type RepoBytes struct {
	Cfg    config.Config
	Logg   logg.Logger
	Store  db.DataBase[*sql.DB]
	SqlDB  *sql.DB
	NameTb string
}

func NewRepoBytes(cfg config.Config, logger logg.Logger, db db.DataBase[*sql.DB]) *RepoBytes {
	return &RepoBytes{
		Cfg:    cfg,
		Logg:   logger,
		Store:  db,
		SqlDB:  db.GetDB(),
		NameTb: "bytes",
	}
}

func (r *RepoBytes) CreateRecord(ctx context.Context, bytes *model.Bytes) (*model.Bytes, error) {
	// TODO
	// Обновить при записи данных в БД:
	// CreatedAt    time.Time
	// UpdatedAt    time.Time

	// *os.File не пишем в БД!

	// bytes.CreatedAt = time.Now()
	// bytes.UpdatedAt = time.Now()

	// r.Store[bytes.Name] = bytes

	return bytes, nil
}

func (r *RepoBytes) ReadRecord(ctx context.Context, bytes *model.Bytes) (*model.Bytes, error) {
	// for k, record := range r.Store {
	// 	if record.Owner == bytes.Owner && k == bytes.Name {
	// 		return record, nil
	// 	}
	// }
	return nil, errs.ErrDataNotFound
}

func (r *RepoBytes) ReadAllRecord(ctx context.Context, bytes *model.Bytes) ([]*model.Bytes, error) {
	// res := make([]*model.Bytes, 0, 10)

	// for _, record := range r.Store {
	// 	if record.Owner == bytes.Owner && record.Type == bytes.Type {
	// 		res = append(res, record)
	// 	}
	// }

	// if len(res) == 0 {
	// 	return nil, errs.ErrDataNotFound
	// }
	return nil, nil
}

func (r *RepoBytes) DeleteRecord(ctx context.Context, text *model.Bytes) (*model.Bytes, error) {
	// for k, record := range r.Store {
	// 	if record.Owner == text.Owner && record.Name == text.Name {
	// 		delete(r.Store, k)
	// 		return record, nil
	// 	}
	// }
	return nil, errs.ErrDataNotFound
}

func (r *RepoBytes) UpdateRecord(ctx context.Context, text *model.Bytes) (*model.Bytes, error) {
	return nil, nil
}
