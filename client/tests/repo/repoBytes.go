package repo

import (
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/errs"
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

	// *os.File не пишем в БД!

	bytes.CreatedAt = time.Now()
	bytes.UpdatedAt = time.Now()

	r.Store[bytes.Name] = bytes

	return bytes, nil
}

func (r *RepoBytes) ReadRecord(bytes *model.Bytes) (*model.Bytes, error) {
	for k, record := range r.Store {
		if record.Owner == bytes.Owner && k == bytes.Name {
			return record, nil
		}
	}
	return nil, errs.ErrDataNotFound
}

func (r *RepoBytes) ReadAllRecord(bytes *model.Bytes) ([]*model.Bytes, error) {
	res := make([]*model.Bytes, 0, 10)

	for _, record := range r.Store {
		if record.Owner == bytes.Owner && record.Type == bytes.Type {
			res = append(res, record)
		}
	}

	if len(res) == 0 {
		return nil, errs.ErrDataNotFound
	}
	return res, nil
}

func (r *RepoBytes) DeleteRecord(text *model.Bytes) (*model.Bytes, error) {
	for k, record := range r.Store {
		if record.Owner == text.Owner && record.Name == text.Name {
			delete(r.Store, k)
			return record, nil
		}
	}
	return nil, errs.ErrDataNotFound
}

func (r *RepoBytes) UpdateRecord(text *model.Bytes) (*model.Bytes, error) {
	return nil, nil
}
