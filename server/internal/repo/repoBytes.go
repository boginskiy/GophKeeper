package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/db"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/infra"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
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
	userID, err := infra.TakeServerValInt64FromCtx(ctx, infra.IDCtx)
	if err != nil {
		return nil, err
	}

	updatedAt := time.Now()

	_, err = r.SqlDB.ExecContext(ctx,
		`INSERT INTO bytes (name, path, sent_size, received_size, type, updated_at, created_at, user_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		bytes.Name,
		bytes.Path,
		bytes.SentSize,
		bytes.ReceivedSize,
		bytes.Type,
		utils.ConvertDtStr(updatedAt),
		utils.ConvertDtStr(updatedAt),
		userID)

	if err != nil {
		return nil, err
	}
	bytes.CreatedAt = updatedAt
	bytes.UpdatedAt = updatedAt

	return bytes, nil
}

func (r *RepoBytes) ReadRecord(ctx context.Context, bytes *model.Bytes) (*model.Bytes, error) {
	userID, err := infra.TakeServerValInt64FromCtx(ctx, infra.IDCtx)
	if err != nil {
		return nil, err
	}
	row := r.SqlDB.QueryRowContext(ctx,
		`SELECT name, received_size, type, updated_at
		 FROM bytes 
		 WHERE user_id = $1 AND name = $2`,
		userID, bytes.Name)

	err = row.Scan(
		&bytes.Name,
		&bytes.ReceivedSize,
		&bytes.Type,
		&bytes.UpdatedAt,
	)

	if err == sql.ErrNoRows { // Нет данных
		return nil, errs.ErrDataNotFound
	}
	if err != nil { // Иные ошибки
		return nil, err
	}

	return bytes, nil
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

// CREATE TABLE bytes (
//     id SERIAL PRIMARY KEY,
//     name TEXT NOT NULL,
//     path TEXT NOT NULL,
//     sent_size BIGINT,
//     received_size BIGINT,
//     type VARCHAR(15) NOT NULL,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
// );

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
