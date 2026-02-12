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
	Logger logg.Logger
	Store  db.DataBase[*sql.DB]
	SqlDB  *sql.DB
	NameTb string
}

func NewRepoBytes(cfg config.Config, logger logg.Logger, db db.DataBase[*sql.DB]) *RepoBytes {
	return &RepoBytes{
		Cfg:    cfg,
		Logger: logger,
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

	currentTime := time.Now()

	_, err = r.SqlDB.ExecContext(ctx,
		`INSERT INTO bytes (name, path, sent_size, received_size, type, updated_at, created_at, user_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		bytes.Name,
		bytes.Path,
		bytes.SentSize,
		bytes.ReceivedSize,
		bytes.Type,
		currentTime,
		currentTime,
		userID)

	if err != nil {
		return nil, err
	}
	bytes.CreatedAt = utils.ConversDtToTableView(currentTime)
	bytes.UpdatedAt = utils.ConversDtToTableView(currentTime)

	return bytes, nil
}

func (r *RepoBytes) ReadRecord(ctx context.Context, bytes *model.Bytes) (*model.Bytes, error) {
	userID, err := infra.TakeServerValInt64FromCtx(ctx, infra.IDCtx)
	if err != nil {
		return nil, err
	}
	row := r.SqlDB.QueryRowContext(ctx,
		`SELECT name, path, received_size, type, created_at
		 FROM bytes 
		 WHERE user_id = $1 AND name = $2`,
		userID, bytes.Name)

	err = row.Scan(
		&bytes.Name,
		&bytes.Path,
		&bytes.ReceivedSize,
		&bytes.Type,
		&bytes.CreatedAt,
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
	userID, err := infra.TakeServerValInt64FromCtx(ctx, infra.IDCtx)
	if err != nil {
		return nil, err
	}

	rows, err := r.SqlDB.QueryContext(ctx,
		`SELECT name, received_size, updated_at 
 		 FROM bytes
		 WHERE user_id = $1 AND type = $2;`,
		userID, bytes.Type)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Читаем данные
	records := make([]*model.Bytes, 0, 10)

	for rows.Next() {
		record := model.Bytes{}
		err := rows.Scan(&record.Name, &record.ReceivedSize, &record.UpdatedAt)
		if err != nil {
			continue
		}
		records = append(records, &record)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if len(records) == 0 {
		return nil, errs.ErrDataNotFound
	}
	return records, nil
}

func (r *RepoBytes) DeleteRecord(ctx context.Context, bytes *model.Bytes) (*model.Bytes, error) {
	userID, err := infra.TakeServerValInt64FromCtx(ctx, infra.IDCtx)
	if err != nil {
		return nil, err
	}

	deleteTime := time.Now()

	result, err := r.SqlDB.ExecContext(ctx,
		`DELETE
		 FROM bytes
		 WHERE user_id = $1 AND name = $2;`,
		userID, bytes.Name)

	// Проверка количеств затронутых изменений.
	cntChange, err := result.RowsAffected()
	if err != nil && cntChange == 0 {
		return nil, errs.ErrDataNotFound
	}

	bytes.UpdatedAt = utils.ConversDtToTableView(deleteTime)
	return bytes, nil
}

func (r *RepoBytes) UpdateRecord(ctx context.Context, text *model.Bytes) (*model.Bytes, error) {
	return nil, nil
}
