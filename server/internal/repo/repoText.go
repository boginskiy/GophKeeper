package repo

import (
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/model"
)

type RepoText struct {
	Store map[string]*model.Text
}

func NewRepoText() *RepoText {
	return &RepoText{Store: make(map[string]*model.Text, 10)}
}

func (r *RepoText) CreateRecord(text *model.Text) (*model.Text, error) {
	// TODO
	// Обновить при записи данных в БД:
	// CreatedAt    time.Time
	// UpdatedAt    time.Time

	text.CreatedAt = time.Now()
	text.UpdatedAt = time.Now()

	r.Store[text.Name] = text

	return text, nil
}

func (r *RepoText) ReadRecord(text *model.Text) (*model.Text, error) {
	for k, record := range r.Store {
		if record.Owner == text.Owner && k == text.Name {
			return record, nil
		}
	}
	return nil, errs.ErrDataNotFound
}

func (r *RepoText) ReadAllRecord(text *model.Text) ([]*model.Text, error) {
	res := make([]*model.Text, 0, 10)

	for _, record := range r.Store {
		if record.Owner == text.Owner && record.Type == text.Type {
			res = append(res, record)
		}
	}

	if len(res) == 0 {
		return nil, errs.ErrDataNotFound
	}
	return res, nil
}

func (r *RepoText) UpdateRecord(text *model.Text) (*model.Text, error) {
	for _, record := range r.Store {

		if record.Owner == text.Owner && record.Name == text.Name {
			record.UpdatedAt = time.Now()
			record.Tx = text.Tx
			return record, nil
		}
	}
	return nil, errs.ErrDataNotFound
}

func (r *RepoText) DeleteRecord(text *model.Text) (*model.Text, error) {
	for k, record := range r.Store {
		if record.Owner == text.Owner && record.Name == text.Name {
			delete(r.Store, k)
			return record, nil
		}
	}
	return nil, errs.ErrDataNotFound
}
