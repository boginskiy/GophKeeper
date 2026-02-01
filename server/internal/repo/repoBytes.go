package repo

import (
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/model"
)

type RepoBytes struct {
	Store map[string]*model.Bytes
}

func NewRepoBytes() *RepoBytes {
	return &RepoBytes{Store: make(map[string]*model.Bytes, 10)}
}

func (r *RepoBytes) CreateRecord(bytes *model.Bytes) (*model.Bytes, error) {
	// TODO
	// Обновить при записи данных в БД:
	// CreatedAt    time.Time
	// UpdatedAt    time.Time

	bytes.CreatedAt = time.Now()
	bytes.UpdatedAt = time.Now()

	r.Store[bytes.Path] = bytes

	return bytes, nil
}
