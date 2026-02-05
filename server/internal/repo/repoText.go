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

type RepoText struct {
	Cfg    config.Config
	Logg   logg.Logger
	Store  db.DataBase[*sql.DB]
	SqlDB  *sql.DB
	NameTb string
}

func NewRepoText(cfg config.Config, logger logg.Logger, db db.DataBase[*sql.DB]) *RepoText {
	return &RepoText{
		Cfg:    cfg,
		Logg:   logger,
		Store:  db,
		SqlDB:  db.GetDB(),
		NameTb: "texts",
	}
}

func (r *RepoText) CreateRecord(ctx context.Context, text *model.Text) (*model.Text, error) {
	userID, err := infra.TakeServerValInt64FromCtx(ctx, infra.IDCtx)
	if err != nil {
		return nil, err
	}

	updatedAt := time.Now()

	_, err = r.SqlDB.ExecContext(ctx,
		`INSERT INTO texts (name, type, content, created_at, updated_at, user_id)
		 VALUES ($1, $2, $3, $4, $5, $6);`,
		text.Name, text.Type, text.Content, utils.ConvertDtStr(updatedAt), utils.ConvertDtStr(updatedAt), userID)

	if err != nil {
		return nil, err
	}

	text.CreatedAt = updatedAt
	text.UpdatedAt = updatedAt
	return text, nil

}

func (r *RepoText) ReadRecord(ctx context.Context, text *model.Text) (*model.Text, error) {
	userID, err := infra.TakeServerValInt64FromCtx(ctx, infra.IDCtx)
	if err != nil {
		return nil, err
	}
	row := r.SqlDB.QueryRowContext(ctx,
		`SELECT type, content, created_at, updated_at
		 FROM texts 
		 WHERE user_id = $1 AND name = $2;`,
		userID, text.Name)

	if err := row.Scan(
		&text.Type,
		&text.Content,
		&text.CreatedAt,
		&text.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return text, nil
}

func (r *RepoText) ReadAllRecord(ctx context.Context, text *model.Text) ([]*model.Text, error) {
	userID, err := infra.TakeServerValInt64FromCtx(ctx, infra.IDCtx)
	if err != nil {
		return nil, err
	}
	rows, err := r.SqlDB.QueryContext(ctx,
		`SELECT name, type, content, created_at, updated_at 
 		 FROM texts 
		 WHERE user_id = $1 AND type = $2;`,
		userID, text.Type)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Читаем данные
	records := make([]*model.Text, 0, 10)

	for rows.Next() {
		record := model.Text{}
		err := rows.Scan(&record.Name, &record.Type, &record.Content, &record.CreatedAt, &record.UpdatedAt)
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

func (r *RepoText) UpdateRecord(ctx context.Context, text *model.Text) (*model.Text, error) {
	userID, err := infra.TakeServerValInt64FromCtx(ctx, infra.IDCtx)
	if err != nil {
		return nil, err
	}

	updatedAt := time.Now()

	result, err := r.SqlDB.ExecContext(ctx,
		`UPDATE texts
		 SET content = $1, updated_at = $2
		 WHERE user_id = $3 AND name = $4;`,
		text.Content, utils.ConvertDtStr(updatedAt), userID, text.Name)

	// Проверка количеств затронутых изменений.
	cntChange, err := result.RowsAffected()
	if err != nil && cntChange == 0 {
		return nil, errs.ErrDataNotFound
	}

	text.UpdatedAt = updatedAt
	return text, nil
}

func (r *RepoText) DeleteRecord(ctx context.Context, text *model.Text) (*model.Text, error) {
	userID, err := infra.TakeServerValInt64FromCtx(ctx, infra.IDCtx)
	if err != nil {
		return nil, err
	}

	deletedAt := time.Now()

	result, err := r.SqlDB.ExecContext(ctx,
		`DELETE
		 FROM texts
		 WHERE user_id = $1 AND name = $2;`,
		userID, text.Name)

	// Проверка количеств затронутых изменений.
	cntChange, err := result.RowsAffected()
	if err != nil && cntChange == 0 {
		return nil, errs.ErrDataNotFound
	}

	text.UpdatedAt = deletedAt
	return text, nil
}
